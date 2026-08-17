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

