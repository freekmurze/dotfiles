# Authentication

## Login

Run `mailcoach login` to authenticate. The command walks through an interactive flow:

1. **Choose instance type**: Mailcoach Cloud (team name) or self-hosted (full URL)
2. **Create API token**: The CLI offers to open your browser to the token creation page
3. **Enter token**: Paste your API token (input is hidden)
4. **Verification**: The CLI verifies the token against the API

```bash
mailcoach login
```

Credentials are stored in `~/.mailcoach/config.json` containing the API token and instance base URL.

## Verify authentication

Run any list command. A successful response means authentication is working:

```bash
mailcoach list-email-lists
```

If you get a 401 error, the token is invalid or expired. Run `mailcoach login` again.

## Logout

```bash
mailcoach logout
```

This removes the stored credentials from `~/.mailcoach/config.json`.

## Important

- Each Mailcoach instance has a unique URL. The CLI stores this during login.
- Mailcoach Cloud URLs follow the pattern `https://<team>.mailcoach.app`
- Self-hosted instances can be any URL
