# laravel-lsp

Wires the first-party [Laravel language server](https://packagist.org/packages/laravel/lsp)
into Claude Code, which has no official plugin for it.

It teaches the editor about the application rather than the language: config keys,
route names, view paths, translation strings, middleware aliases, and container
bindings, across both `.php` and `.blade.php` files.

Complementary to `php-lsp` (Intelephense). Intelephense answers "does this method
exist on this type", Laravel LSP answers "does this route name resolve".

## Requires

```sh
composer global require laravel/lsp
```

Composer's global bin directory must be on `PATH`. `phpEnvironment` defaults to
`auto`, which finds Herd and Valet without configuration.

## Why this only claims `.blade.php`

Claude Code registers one LSP server per file extension. `php-lsp` (Intelephense)
claims `.php`, so this plugin takes `.blade.php`, which nothing else covers.

That splits the work sensibly: Intelephense gives type errors, undefined methods,
and references in PHP, and Laravel LSP gives route names, view paths, config keys,
and translation strings in Blade.

The cost is that Laravel's framework awareness does not apply inside `.php` files,
so a wrong `route()` or `config()` key in a controller is not flagged. To swap the
priority, disable `php-lsp` and add `".php": "php"` back to `.lsp.json`.

Laravel LSP is not a substitute for Intelephense on `.php`. Asked over an LSP
`initialize` handshake (v0.0.31), it advertises only:

```
codeActionProvider   quickfix
completionProvider   triggers " ' . | x - : @
definitionProvider   true
documentLinkProvider true
hoverProvider        true
```

There is no `referencesProvider`, no document or workspace symbol provider, no
signature help, no rename, and no type diagnostics. The completion triggers are
string and Blade characters, with no `$` or `>`, so it does not complete member
access. It knows the application, not the language.
