---
name: spatie-package-skeleton
description: Guide for creating PHP and Laravel packages using Spatie's package-skeleton-laravel and package-skeleton-php templates. Use when the user wants to create a new PHP or Laravel package, scaffold a package. Also use when building customizable packages — covers proven patterns for extensibility (events, configurable models/jobs, action classes) instead of config option creep.
metadata:
  author: Spatie
  tags:
    - php
    - laravel
    - package
    - skeleton
    - open-source
---

# Creating a Laravel Package with Spatie's Skeleton

## Prerequisites

- `gh` CLI installed and authenticated
- `php` available in PATH
- `composer` available in PATH

## Workflow

### 1. Gather Package Details

Ask the user for:
- **Vendor name** (e.g. `spatie`) — the GitHub org or username
- **Package name** (e.g. `laravel-cool-feature`) — the repo/package name
- **Package description** — one-liner for composer.json
- **Visibility** — public or private (default: public)

Use defaults where sensible:
- Author name: from `git config user.name`
- Author email: from `git config user.email`
- Author username: from `gh auth status`
- Vendor namespace: PascalCase of vendor name (e.g. `Spatie`)
- Class name: TitleCase of package name without `laravel-` prefix (e.g. `CoolFeature`)

### 2. Create the Repository from Template

```bash
gh repo create <vendor>/<package-name> --template spatie/package-skeleton-laravel --public --clone
cd <package-name>
```

If the user wants a private repo, use `--private` instead of `--public`.

### 3. Configure the Package (Manual Replacement)

**WARNING**: Do NOT pipe stdin to `configure.php`. The script's child processes (`gh auth status`, `git log`, `git config`) consume lines from the piped stdin, causing inputs to shift and produce garbled results. Instead, do the replacements manually:

1. Run `sed` to replace all placeholder strings across the repo:

```bash
find . -type f -not -path './.git/*' -not -path './vendor/*' -not -name 'configure.php' -exec sed -i '' \
  -e 's/:author_name/Author Name/g' \
  -e 's/:author_username/authorusername/g' \
  -e 's/author@domain\.com/author@email.com/g' \
  -e 's/:vendor_name/Vendor Name/g' \
  -e 's/:vendor_slug/vendorslug/g' \
  -e 's/VendorName/VendorNamespace/g' \
  -e 's/:package_slug_without_prefix/package-without-prefix/g' \
  -e 's/:package_slug/package-name/g' \
  -e 's/:package_name/package-name/g' \
  -e 's/:package_description/Package description here/g' \
  -e 's/Skeleton/ClassName/g' \
  -e 's/skeleton/package-name/g' \
  -e 's/migration_table_name/package_without_prefix/g' \
  -e 's/variable/variableName/g' \
  {} +
```

**Important**: The order of `-e` flags matters. Replace `:package_slug_without_prefix` before `:package_slug` to avoid partial matches. Replace `Skeleton` (PascalCase) before `skeleton` (lowercase).

2. Rename the skeleton files:

```bash
mv src/Skeleton.php src/ClassName.php
mv src/SkeletonServiceProvider.php src/ClassNameServiceProvider.php
mv src/Facades/Skeleton.php src/Facades/ClassName.php
mv src/Commands/SkeletonCommand.php src/Commands/ClassNameCommand.php
mv config/skeleton.php config/package-without-prefix.php
mv database/migrations/create_skeleton_table.php.stub database/migrations/create_package_without_prefix_table.php.stub
```

3. Delete `configure.php` and run `composer install`:

```bash
rm configure.php
composer install
```

Use a longer timeout (5 minutes) for `composer install`.

### 4. Verify Setup

After the script completes:

```bash
# Check the directory structure
ls -la src/
# Verify composer.json looks correct
cat composer.json | head -20
# Check tests passed during setup
```

### 5. Initial Commit and Push

The configure script modifies all files but doesn't commit. Create the initial commit:

```bash
git add -A
git commit -m "Configure package skeleton"
git push -u origin main
```

### 6. Report to User

Tell the user:
- The repo URL (e.g. `https://github.com/<vendor>/<package-name>`)
- The namespace (e.g. `VendorNamespace\ClassName`)
- Key files to start editing:
  - `src/<ClassName>.php` — main package class
  - `src/<ClassName>ServiceProvider.php` — service provider
  - `config/<package-slug>.php` — configuration
  - `tests/` — test directory

## Post-Setup Reference

### Directory Structure

```
src/
  YourClass.php                    # Main package class
  YourClassServiceProvider.php     # Service provider (uses spatie/laravel-package-tools)
  Facades/YourClass.php            # Facade
  Commands/YourClassCommand.php    # Artisan command stub
config/
  your-package.php                 # Published config file
database/
  factories/ModelFactory.php       # Factory template (commented out)
  migrations/create_table.php.stub # Migration stub
resources/views/                   # Blade views
tests/
  TestCase.php                     # Extends Orchestra\Testbench\TestCase
  ArchTest.php                     # Architecture tests (no dd/dump/ray)
  ExampleTest.php                  # Starter test
  Pest.php                         # Pest config binding TestCase
```

### Service Provider Configuration

Uses `spatie/laravel-package-tools`:

```php
public function configurePackage(Package $package): void
{
    $package
        ->name('your-package')
        ->hasConfigFile()
        ->hasViews()
        ->hasMigration('create_your_package_table')
        ->hasCommand(YourClassCommand::class);
}
```

Remove methods you don't need. Delete corresponding directories/files too:

- No database? Delete `database/` and remove `->hasMigration()`
- No commands? Delete `src/Commands/` and remove `->hasCommand()`
- No views? Delete `resources/views/` and remove `->hasViews()`
- No facade? Delete `src/Facades/` and remove facade alias from `composer.json` `extra.laravel.aliases`
- No config? Delete `config/` and remove `->hasConfigFile()`

### Testing

```bash
composer test       # Run tests
composer format     # Run code style fixer
composer analyse    # Run static analysis
```

### Adding an Install Command

```php
use Spatie\LaravelPackageTools\Commands\InstallCommand;

$package->hasInstallCommand(function (InstallCommand $command) {
    $command
        ->publishConfigFile()
        ->publishMigrations()
        ->askToRunMigrations()
        ->askToStarRepoOnGitHub('vendor/package-name');
});
```

## API Design Principles

- **Optimize for easy usage.** The API exposed to users should be as simple as possible. Every public method, facade call, and middleware should feel obvious and require minimal setup.
- **Use well-named methods.** Method names should be intuitive and self-documenting. Prefer descriptive names over terse ones — the user should understand what a method does without reading its implementation. Use verb-first method names (`clear()`, `forget()`, `save()`).
- **Follow Spatie PHP/Laravel guidelines.** All code must follow the conventions described in the `spatie-guidelines` skill.

## Package patters

### Fluent/Chainable APIs

Builder-style classes where every setter returns `$this`. Users should be able to chain configuration calls naturally.

```php
Pdf::view('invoice', $data)->format('a4')->landscape()->save('invoice.pdf');
```

### Sensible Defaults

The package should work well out of the box with zero configuration. Only require explicit setup for non-standard use cases. Provide safe defaults in the config file and apply them when values aren't explicitly set.

### Facade + Factory for Clean State

Back facades with a factory that creates a fresh builder per call to prevent state bleed between requests.

```php
// Factory intercepts calls via __call() to create fresh builder instances
class PdfFactory {
    public function __call($method, $parameters) {
        return (clone $this->builder)->$method(...$parameters);
    }
}
```

### Enums Over Strings

Use PHP enums for any fixed set of options instead of string constants. This gives type safety and IDE support.

### Value Objects for Options

Group related settings into small readonly classes (like `PdfOptions`, `ScreenshotOptions`) rather than passing many loose parameters between layers.

### Descriptive Exception Classes

Name exceptions after what went wrong and provide static factory methods for specific scenarios with helpful error messages:

```php
class CouldNotGeneratePdf extends Exception
{
    public static function browsershotNotInstalled(): static
    {
        return new static('To use Browsershot, install it via `composer require spatie/browsershot`.');
    }
}
```

### Traits for Cross-Cutting Concerns

Use `Conditionable` (for `when()`/`unless()` chaining), `Macroable` (for runtime extension), and `Dumpable` (for debugging) on builder classes.

### Small Interfaces for Extensibility

Define interfaces for components users might want to swap. Keep them small — one or two methods is ideal:

```php
interface PdfDriver {
    public function generatePdf(string $html, ...): string;
    public function savePdf(string $html, ..., string $path): void;
}
```

### Config-Driven Class Bindings

Let users swap implementations via config rather than requiring service provider overrides:

```php
// config/your-package.php
'driver' => env('LARAVEL_PDF_DRIVER', 'browsershot'),
'cache_profile' => App\CacheProfiles\CustomCacheProfile::class,
'hasher' => App\Hashers\CustomHasher::class,
```

### Testing Fakes with Rich Assertions

Provide a `::fake()` method on the facade that swaps in a fake builder. Track calls and offer assertion methods:

```php
Pdf::fake();
// ... code that generates PDFs ...
Pdf::assertSaved(fn ($pdf, $path) => $path === 'invoice.pdf');
Pdf::assertQueued();
Pdf::assertNotQueued();
```

### Events at Key Moments

Fire events for important lifecycle moments so users can hook into the workflow without modifying package code.

### Anti-Pattern: Config Option Creep

Don't add small config options for every customization request. Instead, give users full control via class extension.

### Pattern: Events Instead of Hook Config Options

Fire events and let users listen:

```php
event(new TransformerStarting($transformer, $url));
$transformer->transform();
event(new TransformerEnded($transformer, $url, $result));
```

### Pattern: Configurable Models

Let users specify their own model class in config:

```php
// config
'model' => Spatie\Package\Models\Result::class,

// In package code — always resolve from config:
$model = config('your-package.model');
$model::find($id);
```

### Pattern: Configurable Jobs

Let users specify their own job class in config:

```php
'process_job' => Spatie\Package\Jobs\ProcessJob::class,
```

### Pattern: Action Classes

Wrap small pieces of functionality in action classes registered in config:

```php
'actions' => [
    'fetch_content' => Spatie\Package\Actions\FetchContentAction::class,
],
```

Users override by extending and registering their custom action.

### Queued Operations with Callbacks

For expensive operations, provide `saveQueued()` that returns a wrapper around `PendingDispatch` with `then()`/`catch()` callbacks:

```php
Pdf::view('invoice', $data)
    ->saveQueued('invoice.pdf')
    ->then(fn ($path) => /* success */)
    ->catch(fn ($e) => /* failure */)
    ->onQueue('pdfs');
```

### Consistent Naming Conventions

- Suffix event classes with `Event`
- Suffix notification classes with `Notification`
- Suffix config data classes with `Config`
- Use `{Service}Driver` for driver implementations
- Use `Could Not...` for exception classes
- Use `Fake` prefix for test doubles

### Config File Comments

Always add block comments above each config key or group explaining what it does:

```php
return [
    /*
     * When disabled, the middleware will not convert any responses.
     */
    'enabled' => env('PACKAGE_ENABLED', true),

    /*
     * The driver used to perform the operation.
     * Supported: "local", "cloud"
     */
    'driver' => env('PACKAGE_DRIVER', 'local'),

    'cache' => [
        /*
         * How long results should be cached, in seconds.
         */
        'ttl' => (int) env('PACKAGE_CACHE_TTL', 3600),
    ],
];
```

Use `/* */` block comments (not `//`). Mention supported values, defaults, and any non-obvious behavior. Keep comments concise — one to three lines.

### Miscellaneous

- do not add `down` methods to migration
- do not use else statements — return early instead
- do not use compound if statements — split into multiple ifs or use guard clauses
- use `protected` visibility over `private`