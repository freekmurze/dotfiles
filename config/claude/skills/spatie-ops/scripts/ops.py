#!/usr/bin/env python3
"""Provider API requests with private credential references and redacted output."""

import argparse
import json
import os
import re
import stat
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

CONFIG_PATH = Path.home() / '.dotfiles-custom/spatie-ops/credentials.json'
ALIASES = {'do': 'digitalocean', 'digital-ocean': 'digitalocean', 'laravel-cloud': 'cloud'}
BASES = {
    'digitalocean': 'https://api.digitalocean.com/v2',
    'forge': 'https://forge.laravel.com/api',
    'cloud': 'https://cloud.laravel.com/api',
    'openprovider': 'https://api.openprovider.eu/v1beta',
}
REDACTED = '[redacted]'


class OpsError(Exception):
    pass


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, request, response, code, message, headers, new_url):
        return None


def read_config(path):
    path = Path(path).expanduser().resolve()
    try:
        info = path.stat()
    except FileNotFoundError:
        raise OpsError('Private configuration is missing. See spatie-ops/references/authentication.md.') from None
    if not stat.S_ISREG(info.st_mode) or info.st_uid != os.getuid() or stat.S_IMODE(info.st_mode) != 0o600:
        raise OpsError('Private configuration must be a user-owned regular file with mode 600.')
    try:
        config = json.loads(path.read_text())
    except (ValueError, OSError):
        raise OpsError('Private configuration could not be read as JSON.') from None
    if not isinstance(config, dict) or not isinstance(config.get('providers'), dict):
        raise OpsError('Private configuration must contain a providers object.')
    return config


class Client:
    def __init__(self, config):
        self.providers = config['providers']
        self.tokens = {}
        self.secrets = set()
        self.opener = urllib.request.build_opener(NoRedirect)

    @classmethod
    def from_file(cls, path=None):
        return cls(read_config(path or os.environ.get('SPATIE_OPS_CONFIG', CONFIG_PATH)))

    def op(self, arguments):
        try:
            result = subprocess.run(['op', *arguments], capture_output=True, text=True, timeout=55)
        except subprocess.TimeoutExpired:
            raise OpsError('1Password authorization timed out. Unlock it and authorize the request.') from None
        except FileNotFoundError:
            raise OpsError('Install the 1Password op CLI before making this request.') from None
        if result.returncode:
            raise OpsError('1Password could not provide the credential. Check account access and authorization.')
        return result.stdout.strip()

    def remember_secret(self, value):
        if not isinstance(value, str) or not value:
            raise OpsError('Credential is missing or empty.')
        self.secrets.add(value)
        return value

    def read_ref(self, reference):
        if not isinstance(reference, str) or not reference.startswith('op://'):
            raise OpsError('Configure a 1Password secret reference for this provider.')
        return self.remember_secret(self.op(['read', reference]))

    def http(self, url, method, token=None, payload=None):
        headers = {'Accept': 'application/json', 'User-Agent': 'spatie-ops/1.0'}
        if token:
            headers['Authorization'] = 'Bearer ' + token
        body = None
        if payload is not None:
            self.remember_payload_secrets(payload)
            headers['Content-Type'] = 'application/json'
            body = json.dumps(payload).encode()
        request = urllib.request.Request(url, data=body, headers=headers, method=method)
        try:
            response = self.opener.open(request, timeout=30)
        except urllib.error.HTTPError as error:
            response = error
        except (urllib.error.URLError, TimeoutError, OSError):
            raise OpsError('HTTPS request failed. Check connectivity and provider availability.') from None
        with response:
            content = response.read()
        try:
            data = json.loads(content) if content else None
        except (ValueError, UnicodeDecodeError):
            data = {'message': 'Non-JSON response omitted.'}
        return response.status, data

    def remember_payload_secrets(self, value, environment=False):
        if isinstance(value, list):
            for item in value:
                self.remember_payload_secrets(item, environment)
        if isinstance(value, dict):
            environment_value = ('key' in value or environment) and 'value' in value
            for key, item in value.items():
                if isinstance(item, str) and item:
                    if self.sensitive_key(key) or (environment_value and key == 'value'):
                        self.secrets.add(item)
                self.remember_payload_secrets(item, key in {'environment_variables', 'variables'})

    def token(self, provider):
        if provider in self.tokens:
            return self.tokens[provider]
        settings = self.providers[provider]
        token = None
        if provider == 'cloud' and settings.get('cli_config'):
            try:
                saved = json.loads(Path(settings['cli_config']).expanduser().read_text())
                tokens = saved.get('api_tokens', [])
                token = tokens[0] if tokens else None
            except (OSError, ValueError):
                pass
        if provider == 'openprovider':
            if settings.get('login_item'):
                fields = json.loads(self.op([
                    'item', 'get', settings['login_item'], '--fields',
                    'label=username,label=password', '--reveal', '--format', 'json',
                ]))
                if isinstance(fields, dict):
                    fields = fields.get('fields', [fields])
                values = {field.get('label', field.get('id')): field.get('value') for field in fields}
                username = self.remember_secret(values.get('username'))
                password = self.remember_secret(values.get('password'))
            else:
                username = self.read_ref(settings.get('username_ref'))
                password = self.read_ref(settings.get('password_ref'))
            status, data = self.http(BASES[provider] + '/auth/login', 'POST', payload={
                'username': username, 'password': password,
            })
            if status != 200 or not isinstance(data, dict) or not data.get('data', {}).get('token'):
                raise OpsError('Openprovider authentication failed. Verify the configured API login.')
            token = data['data']['token']
        if not token:
            token = self.read_ref(settings.get('token_ref'))
        token = self.remember_secret(token)
        self.tokens[provider] = token
        return token

    def url(self, provider, path):
        if provider not in BASES or provider not in self.providers:
            raise OpsError('Provider is unsupported or not configured.')
        parsed = urllib.parse.urlsplit(path)
        decoded = parsed.path
        for _ in range(5):
            previous = decoded
            decoded = urllib.parse.unquote(decoded)
            if decoded == previous:
                break
        if parsed.scheme or parsed.netloc or parsed.fragment or not path.startswith('/'):
            raise OpsError('Use a provider-relative path beginning with a single slash.')
        if decoded.startswith('//') or '\\' in decoded or '..' in decoded.split('/'):
            raise OpsError('Path must stay inside the configured provider API.')
        for name, _ in urllib.parse.parse_qsl(parsed.query):
            if self.sensitive_key(name):
                raise OpsError('Send credentials in a JSON body, not URL query parameters.')
        base = BASES[provider]
        if provider == 'forge':
            organization = self.providers[provider].get('organization')
            if not organization:
                raise OpsError('Configure the Forge organization in the private mapping.')
            base += '/orgs/' + urllib.parse.quote(organization, safe='')
        return base + path

    def request(self, provider, method, path, payload=None):
        provider = ALIASES.get(provider, provider)
        method = method.upper()
        if method not in {'GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD'}:
            raise OpsError('Unsupported HTTP method.')
        url = self.url(provider, path)
        return self.http(url, method, self.token(provider), payload)

    @staticmethod
    def sensitive_key(key):
        key = re.sub(r'[^a-z0-9]', '', key.lower())
        return key == 'key' or any(part in key for part in (
            'password', 'passwd', 'secret', 'credential', 'token', 'authorization',
            'cookie', 'apikey', 'privatekey', 'accesskey', 'signature',
        ))

    def redact(self, value, environment=False):
        if isinstance(value, dict):
            environment_value = ('key' in value or environment) and 'value' in value
            return {
                key: REDACTED if (
                    self.sensitive_key(key) and not (key == 'key' and environment_value)
                ) or (key == 'value' and environment_value) else self.redact(item, key in {'environment_variables', 'variables'})
                for key, item in value.items()
            }
        if isinstance(value, list):
            return [self.redact(item, environment) for item in value]
        if isinstance(value, str):
            for secret in sorted(self.secrets, key=len, reverse=True):
                value = value.replace(secret, REDACTED)
            value = re.sub(r'(?:gh[opusr]_|github_pat_)[A-Za-z0-9_]+', REDACTED, value)
            value = re.sub(r'eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+', REDACTED, value)
            value = re.sub(r'op://[^\s"<>]+', REDACTED, value)
            value = re.sub(r'(https?://[^\s?]+)\?[^\s]*X-Amz-[^\s]*', r'\1?[redacted]', value)
            return value
        return value


def private_destination(path):
    destination = Path(path).expanduser()
    destination.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
    for directory in [destination.parent.resolve(), *destination.parent.resolve().parents]:
        if (directory / '.git').exists():
            raise OpsError('Private output must be outside a Git working tree.')
    git = subprocess.run(['git', '-C', str(destination.parent.resolve()), 'rev-parse', '--show-toplevel'], capture_output=True)
    if git.returncode == 0:
        raise OpsError('Private output must be outside a Git working tree.')
    if destination.is_symlink():
        raise OpsError('Private output cannot be a symlink.')
    if destination.exists():
        info = destination.stat()
        if not stat.S_ISREG(info.st_mode) or info.st_uid != os.getuid() or info.st_nlink != 1:
            raise OpsError('Private output must be a user-owned regular file without hard links.')
    return destination


def private_output(path, data):
    destination = private_destination(path)
    try:
        fd = os.open(destination, os.O_WRONLY | os.O_CREAT | os.O_NOFOLLOW, 0o600)
    except OSError:
        raise OpsError('Private output could not be opened safely.') from None
    with os.fdopen(fd, 'w') as stream:
        info = os.fstat(stream.fileno())
        if not stat.S_ISREG(info.st_mode) or info.st_uid != os.getuid() or info.st_nlink != 1:
            raise OpsError('Private output must be a user-owned regular file without hard links.')
        os.fchmod(stream.fileno(), 0o600)
        stream.truncate(0)
        json.dump(data, stream, indent=2)
        stream.write('\n')


def selected_fields(data, fields):
    result = {}
    for field in fields.split(','):
        item = data
        for part in field.split('.'):
            item = item.get(part) if isinstance(item, dict) else None
        result[field] = item
    return result


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--config', help='Private credential mapping (default: SPATIE_OPS_CONFIG or ~/.dotfiles-custom/spatie-ops/credentials.json)')
    commands = parser.add_subparsers(dest='command', required=True)
    commands.add_parser('status', help='Check private configuration without accessing credentials')
    request = commands.add_parser('request', help='Make an authenticated provider API request')
    request.add_argument('provider', choices=[*BASES, *ALIASES])
    request.add_argument('method')
    request.add_argument('path')
    request.add_argument('--input', help='JSON input file, or - for stdin')
    output = request.add_mutually_exclusive_group()
    output.add_argument('--fields', help='Comma-separated dotted paths from the redacted response')
    output.add_argument('--private-output', help='Save the unredacted response to a mode-600 file outside Git')
    args = parser.parse_args(argv)
    try:
        client = Client.from_file(args.config)
        if args.command == 'status':
            print(json.dumps({'configured': [provider for provider in BASES if provider in client.providers]}))
            return 0
        payload = None
        if args.input:
            text = sys.stdin.read() if args.input == '-' else Path(args.input).expanduser().read_text()
            payload = json.loads(text)
        if args.private_output:
            private_destination(args.private_output)
        status, data = client.request(args.provider, args.method, args.path, payload)
        if args.private_output:
            private_output(args.private_output, data)
            print(json.dumps({'status': status, 'private_output_saved': True}))
        else:
            response = client.redact(data)
            if args.fields:
                response = selected_fields(response, args.fields)
            print(json.dumps({'status': status, 'response': response}, indent=2))
        return 0 if 200 <= status < 300 else 1
    except OpsError as error:
        print(json.dumps({'error': str(error)}), file=sys.stderr)
        return 1
    except Exception:
        print(json.dumps({'error': 'Invalid input, configuration or provider response; sensitive details omitted.'}), file=sys.stderr)
        return 1


if __name__ == '__main__':
    raise SystemExit(main())
