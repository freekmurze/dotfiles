import contextlib
import io
import json
import os
import stat
import subprocess
import tempfile
import unittest
import urllib.error
from pathlib import Path
from unittest.mock import Mock, patch

import ops


class OpsTests(unittest.TestCase):
    def setUp(self):
        self.client = ops.Client({'providers': {
            'cloud': {'token_ref': 'op://example/cloud/credential'},
            'digitalocean': {'token_ref': 'op://example/do/credential'},
            'forge': {'token_ref': 'op://example/forge/credential', 'organization': 'example-org'},
            'openprovider': {'username_ref': 'op://example/dns/username', 'password_ref': 'op://example/dns/password'},
        }})

    def test_unknown_secret_values_in_environment_and_bucket_responses_are_redacted(self):
        response = {'data': {'environment_variables': [
            {'key': 'FLARE_KEY', 'value': 'new-reporting-secret'},
            {'key': 'APP_NAME', 'value': 'Example'},
        ], 'bucket': {'access_key_secret': 'new-storage-secret'}, 'apiKey': 'new-api-secret'}}
        output = json.dumps(self.client.redact(response))
        for secret in ['new-reporting-secret', 'Example', 'new-storage-secret', 'new-api-secret']:
            self.assertNotIn(secret, output)
        self.assertIn('FLARE_KEY', output)

    def test_known_credentials_in_unstructured_errors_and_signed_urls_are_redacted(self):
        self.client.remember_secret('fake-live-bearer')
        response = {'message': 'Rejected fake-live-bearer', 'url': 'https://bucket.example/archive.zip?X-Amz-Expires=300&X-Amz-Signature=example'}
        output = json.dumps(self.client.redact(response))
        self.assertNotIn('fake-live-bearer', output)
        self.assertNotIn('X-Amz-Signature', output)

    def test_dns_values_are_preserved_while_named_environment_values_are_redacted(self):
        response = {'data': {'records': [{'name': 'www.example.com', 'value': '203.0.113.10', 'type': 'A'}], 'environment_variables': [{'name': 'APP_KEY', 'value': 'fake-app-key'}]}}
        redacted = self.client.redact(response)
        self.assertEqual(redacted['data']['records'][0]['value'], '203.0.113.10')
        self.assertEqual(redacted['data']['environment_variables'][0]['value'], ops.REDACTED)

    def test_sensitive_payload_values_are_not_echoed_in_errors(self):
        self.client.remember_payload_secrets({'variables': [{'key': 'FLARE_KEY', 'value': 'fake-reporting-key'}]})
        self.assertEqual(self.client.redact('Bad fake-reporting-key'), 'Bad [redacted]')

    def test_host_and_traversal_attempts_are_rejected_before_reading_credentials(self):
        with patch.object(self.client, 'token') as token:
            for path in ['https://evil.example/', '//evil.example/', '/%2f%2fevil.example/', '/../elsewhere', '/%2e%2e/elsewhere', '/%252e%252e/elsewhere', '/\\evil.example/', '/x?token=secret']:
                with self.subTest(path=path), self.assertRaises(ops.OpsError):
                    self.client.request('cloud', 'GET', path)
            token.assert_not_called()

    def test_credentials_are_only_in_request_memory_and_provider_aliases_work(self):
        response = Mock()
        response.status = 200
        response.read.return_value = b'{"droplets":[]}'
        response.__enter__ = Mock(return_value=response)
        response.__exit__ = Mock(return_value=False)
        with patch.object(self.client, 'op', return_value='fake-bearer') as op, patch.object(self.client.opener, 'open', return_value=response) as network:
            self.assertEqual(self.client.request('do', 'GET', '/droplets'), (200, {'droplets': []}))
            op.assert_called_once_with(['read', 'op://example/do/credential'])
            request = network.call_args.args[0]
            self.assertEqual(request.full_url, 'https://api.digitalocean.com/v2/droplets')
            self.assertEqual(request.get_header('Authorization'), 'Bearer fake-bearer')
        self.assertIn(ops.NoRedirect, [type(handler) for handler in self.client.opener.handlers])
        self.assertIsNone(ops.NoRedirect().redirect_request(None, None, 302, '', {}, 'https://evil.example'))

    def test_openprovider_login_is_not_an_authenticated_redirect_or_persistent_token(self):
        with patch.object(self.client, 'read_ref', side_effect=['fake-user', 'fake-password']), patch.object(self.client, 'http', return_value=(200, {'data': {'token': 'fake-dns-bearer'}})) as network:
            self.assertEqual(self.client.token('openprovider'), 'fake-dns-bearer')
            network.assert_called_once_with('https://api.openprovider.eu/v1beta/auth/login', 'POST', payload={'username': 'fake-user', 'password': 'fake-password'})

    def test_cloud_uses_existing_cli_login_without_requesting_1password(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'cloud.json'
            path.write_text(json.dumps({'api_tokens': ['fake-cli-bearer']}))
            self.client.providers['cloud']['cli_config'] = str(path)
            with patch.object(self.client, 'op') as op:
                self.assertEqual(self.client.token('cloud'), 'fake-cli-bearer')
                op.assert_not_called()

    def test_forge_organization_comes_from_private_mapping(self):
        self.assertEqual(self.client.url('forge', '/servers'), 'https://forge.laravel.com/api/orgs/example-org/servers')

    def test_private_configuration_permissions_are_enforced(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'credentials.json'
            path.write_text('{"providers":{}}')
            path.chmod(0o644)
            with self.assertRaises(ops.OpsError):
                ops.read_config(path)
            path.chmod(0o600)
            self.assertEqual(ops.read_config(path), {'providers': {}})

    def test_private_output_is_mode_600_and_rejects_symlinks_and_git_worktrees(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'response.json'
            path.write_text('old public content')
            path.chmod(0o644)
            ops.private_output(path, {'api_key': 'fake-secret'})
            self.assertEqual(stat.S_IMODE(path.stat().st_mode), 0o600)
            self.assertEqual(json.loads(path.read_text()), {'api_key': 'fake-secret'})
            link = Path(directory) / 'link.json'
            link.symlink_to(path)
            with self.assertRaises(ops.OpsError):
                ops.private_output(link, {'new': 'secret'})
            hard_link = Path(directory) / 'hard-link.json'
            os.link(path, hard_link)
            with self.assertRaises(ops.OpsError):
                ops.private_destination(hard_link)
            hard_link.unlink()
            subprocess.run(['git', 'init', '--quiet', directory], check=True)
            with self.assertRaises(ops.OpsError):
                ops.private_output(path, {'new': 'secret'})

    def test_failed_request_is_nonzero_and_error_output_still_redacts(self):
        with patch.object(ops.Client, 'from_file', return_value=self.client), patch.object(self.client, 'request', return_value=(401, {'message': 'Denied', 'token': 'fake-exposed-token'})), contextlib.redirect_stdout(io.StringIO()) as output:
            self.assertEqual(ops.main(['request', 'cloud', 'GET', '/applications']), 1)
            self.assertNotIn('fake-exposed-token', output.getvalue())
            self.assertEqual(json.loads(output.getvalue())['status'], 401)

    def test_private_output_is_rejected_before_an_api_mutation(self):
        with tempfile.TemporaryDirectory() as directory:
            subprocess.run(['git', 'init', '--quiet', directory], check=True)
            with patch.object(ops.Client, 'from_file', return_value=self.client), patch.object(self.client, 'request') as request, contextlib.redirect_stderr(io.StringIO()):
                result = ops.main(['request', 'cloud', 'POST', '/buckets', '--private-output', str(Path(directory) / 'key.json')])
                self.assertEqual(result, 1)
                request.assert_not_called()

    def test_cli_error_does_not_print_subprocess_stderr_or_config_references(self):
        result = subprocess.CompletedProcess(['op'], 1, stdout='fake-leak', stderr='op://private/item/credential fake-leak')
        with patch('ops.subprocess.run', return_value=result):
            with self.assertRaises(ops.OpsError) as error:
                self.client.op(['read', 'op://example/do/credential'])
            self.assertNotIn('fake-leak', str(error.exception))
            self.assertNotIn('op://', str(error.exception))

    def test_status_reads_no_credentials_and_cli_fields_cannot_bypass_redaction(self):
        for argv in [['status'], ['request', 'cloud', 'GET', '/applications', '--fields', 'data.api_key']]:
            with self.subTest(argv=argv), patch.object(ops.Client, 'from_file', return_value=self.client), patch.object(self.client, 'request', return_value=(200, {'data': {'api_key': 'fake-new-key'}})), patch.object(self.client, 'op') as op, contextlib.redirect_stdout(io.StringIO()) as output:
                self.assertEqual(ops.main(argv), 0)
                self.assertNotIn('fake-new-key', output.getvalue())
                self.assertNotIn('op://', output.getvalue())
                op.assert_not_called()


if __name__ == '__main__':
    unittest.main()
