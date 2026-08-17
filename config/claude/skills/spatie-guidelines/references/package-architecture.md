## Spatie Package Architecture

### Package Structure (spatie/package-skeleton-laravel)

```
src/                          # Main source code
src/Commands/                 # Artisan commands
src/Components/               # Blade components
src/Contracts/                # Interfaces
src/Events/                   # Events
src/Exceptions/               # Exception classes
src/Models/                   # Eloquent models
src/Traits/                   # Reusable traits
config/                       # Config file (kebab-case, no `laravel-` prefix)
database/factories/           # Model factories
database/migrations/          # Migrations (stubs)
resources/views/              # Views
routes/                       # Routes
tests/                        # Tests (Pest)
```

### Service Provider Pattern

All Spatie Laravel packages extend `PackageServiceProvider` from `spatie/laravel-package-tools`:

```php
use Spatie\LaravelPackageTools\Package;
use Spatie\LaravelPackageTools\PackageServiceProvider;

class MediaLibraryServiceProvider extends PackageServiceProvider
{
    public function configurePackage(Package $package): void
    {
        $package
            ->name('laravel-medialibrary')     // Package name
            ->hasConfigFile('media-library')   // Config without 'laravel-' prefix
            ->hasMigration('create_media_table')
            ->hasViews('media-library')
            ->hasCommands([
                RegenerateCommand::class,
                ClearCommand::class,
                CleanCommand::class,
            ]);
    }

    public function packageBooted(): void
    {
        // Post-boot logic: observers, event listeners, gate registrations
    }

    public function packageRegistered(): void
    {
        // Bindings, singletons, scoped instances
    }
}
```

**Key lifecycle methods:**

- `configurePackage()` declares assets (config, migrations, views, commands)
- `packageRegistered()` binds interfaces and registers singletons
- `packageBooted()` registers observers, macros, and blade directives

### Namespace Conventions

- Root namespace: `Spatie\PackageName` (for example `Spatie\Permission`, `Spatie\MediaLibrary`)
- Packagist name: `spatie/laravel-package-name` or `spatie/package-name`
- The config file drops the `laravel-` prefix: `spatie/laravel-permission` becomes `config/permission.php`

### Contracts Pattern

Spatie packages use contracts (interfaces) for extensibility:

```php
// src/Contracts/Permission.php
namespace Spatie\Permission\Contracts;

interface Permission
{
    // ...
}

// src/Models/Permission.php, implements the contract
class Permission extends Model implements PermissionContract
{
    // ...
}

// Service provider binds contract to implementation
$this->app->bind(PermissionContract::class,
    fn ($app) => $app->make($app->config['permission.models.permission'])
);
```

Config allows users to swap model implementations:

```php
// config/permission.php
return [
    'models' => [
        'permission' => Spatie\Permission\Models\Permission::class,
        'role' => Spatie\Permission\Models\Role::class,
    ],
    'table_names' => [
        'roles' => 'roles',
        'permissions' => 'permissions',
    ],
];
```

### Model Patterns

- Use `$guarded = []`, not `$fillable`.
- Configurable table names via config: `$this->table = config('permission.table_names.permissions') ?: parent::getTable();`
- Provide static factory methods: `Permission::create()`, `Permission::findByName()`, `Permission::findOrCreate()`
- Throw custom exceptions for domain errors: `PermissionAlreadyExists`, `PermissionDoesNotExist`
- Use traits for shared behavior: `HasRoles`, `InteractsWithMedia`

### Config File Conventions

- Write verbose comments explaining each option in the config file.
- Let users swap class implementations via config (models, generators, and so on).
- Use `snake_case` keys throughout.
- Group related options (`models`, `table_names`, `column_names`).

---

