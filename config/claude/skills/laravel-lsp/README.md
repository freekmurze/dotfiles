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
