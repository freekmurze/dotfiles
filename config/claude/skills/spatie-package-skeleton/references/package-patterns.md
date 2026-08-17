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
