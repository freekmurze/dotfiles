## Testing (Pest)

Spatie packages use **Pest** for testing, with **Orchestra Testbench** for Laravel integration.

### Setup

```php
// tests/Pest.php
use Spatie\YourPackage\Tests\TestCase;

uses(TestCase::class)->in(__DIR__);
```

```php
// tests/TestCase.php
use Orchestra\Testbench\TestCase as Orchestra;

class TestCase extends Orchestra
{
    protected function setUp(): void
    {
        parent::setUp();

        Factory::guessFactoryNamesUsing(
            fn (string $modelName) => 'Spatie\\YourPackage\\Database\\Factories\\' . class_basename($modelName) . 'Factory'
        );
    }

    protected function getPackageProviders($app)
    {
        return [YourPackageServiceProvider::class];
    }

    public function getEnvironmentSetUp($app)
    {
        config()->set('database.default', 'testing');
        // Run migrations, set config overrides
    }
}
```

### Test Style

- Use descriptive test names. Never add comments to tests.
- Follow the arrange, act, assert pattern.

```php
// Use Pest's `it()` syntax with closures
it('can assign a role', function () {
    $user = User::factory()->create();

    $user->assignRole('admin');

    expect($user->hasRole('admin'))->toBeTrue();
});

it('throws when permission does not exist', function () {
    expect(fn () => Permission::findByName('nonexistent'))
        ->toThrow(PermissionDoesNotExist::class);
});

// Chain expectations
it('can create a data object from array', function () {
    $dto = SimpleDto::from('Hello World');

    expect($dto)
        ->toBeInstanceOf(SimpleDto::class)
        ->and($dto->string)->toEqual('Hello World');
});
```

### What to Test

- **Core functionality**: CRUD operations, main feature paths
- **Edge cases**: null inputs, missing data, duplicate entries
- **Custom exceptions**: verify domain errors throw the correct exception types
- **Config overrides**: test that swapping implementations via config works
- **Blade directives and components**: if the package provides them
- **Artisan commands**: test output and side effects

### Test Helpers

- Keep test classes in the same file when possible. Put helper functions and internal test classes *last*, below the `it` and `test` blocks, so the tests read first and the helpers stay support detail:

```php
// At the bottom of the test file, not in a separate file
class ItemAdded extends ShouldBeStored
{
    public function __construct(public string $name) {}
}
```

- For reusable test support classes, put them in `tests/TestSupport/`:

```
tests/
  TestSupport/
    TestModels/
      TestModel.php
      TestModelWithConversion.php
```

### composer.json Testing Stack

```json
{
    "require-dev": {
        "laravel/pint": "^1.14",
        "larastan/larastan": "^3.0",
        "orchestra/testbench": "^10.0||^9.0",
        "pestphp/pest": "^4.0",
        "pestphp/pest-plugin-arch": "^4.0",
        "pestphp/pest-plugin-laravel": "^4.0",
        "phpstan/extension-installer": "^1.4"
    },
    "scripts": {
        "test": "vendor/bin/pest",
        "format": "vendor/bin/pint",
        "analyse": "vendor/bin/phpstan analyse"
    }
}
```

---

