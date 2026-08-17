---
name: spatie-guidelines
description: Spatie's PHP, Laravel, JavaScript and Vue coding guidelines and conventions. Use when writing or reviewing PHP, Laravel, JavaScript, or Vue code for Spatie projects or packages. Covers code style, type declarations, docblocks, control flow, naming, routing, controllers, Blade, validation, comments, testing (Pest), package structure, service providers, and Git workflow. Triggers include "follow Spatie guidelines", "Spatie style", "Spatie package", or any code review for Spatie packages and projects.
license: MIT
metadata:
   author: Spatie
   tags: php, laravel, javascript, best practices, coding standards
---

# Spatie Guidelines

Apply these guidelines when writing code for Spatie projects or contributing to Spatie packages.

**Core principle:** write things the way Laravel intended. If Laravel has a documented way to do something, use it. Deviate only with a clear justification.

---

## PHP Style

### Type System

- Type properties, parameters, and return types. Skip docblocks for fully typed methods.
- Use `?Type` (short nullable), not `Type|null`.
- Always use the `void` return type when a method returns nothing.
- Use constructor property promotion when all properties can be promoted. One per line, trailing comma:

```php
class MyClass {
    public function __construct(
        protected string $firstArgument,
        protected string $secondArgument,
    ) {}
}
```

### Docblocks

- Skip docblocks for fully type-hinted methods unless you need a description.
- Use full sentences with a period for descriptions.
- **Always import classnames, then reference the short name in the docblock.** Never write a fully qualified name inside a docblock:

```php
use Spatie\Url\Url;

/** @return Url */
```

- Use one-line docblocks when possible: `/** @var string */`
- Always add types for iterables, specifying both key and value:

```php
/**
 * @param array<int, MyObject> $myArray
 * @param int $typedArgument
 */
function someFunction(array $myArray, int $typedArgument) {}
```

- Put the most common type first in a multi-type docblock:

```php
/** @var Collection|SomeWeirdVendor\Collection */
```

- Use array shape notation for fixed keys, with each key on its own line:

```php
/** @return array{
    first: SomeClass,
    second: SomeClass
} */
```

- If one parameter needs a docblock, add docblocks for all the other parameters too.

### Code Style

- Follow PSR-1, PSR-2, and PSR-12.
- Use camelCase for non-public-facing strings.
- Don't use `final` by default.
- Prefer string interpolation over concatenation: `"Hi, I am {$name}."`
- Enum values use PascalCase: `case Diamonds;`
- Each trait on its own line with its own `use`:

```php
class MyClass
{
    use TraitA;
    use TraitB;
}
```

- Always import namespaces with `use` statements. Never use inline fully qualified class names such as `\Exception` or `\Illuminate\Support\Facades\Http`.
- Never use single-letter variable names. Write `$exception` instead of `$e`, `$request` instead of `$r`.

### Avoid `private const`

Don't introduce `private const`. Replace each one:

- **Used once:** inline it. When a raw literal would obscure the meaning the name carried, assign it to a descriptively named local variable instead.
- **Used more than once:** turn it into a private property. Use `private static` when the methods using it are static.
- **Referenced from a default parameter value or a PHP attribute**, where no property can be used: inline the literal there.

### Control Flow

- **Happy path last.** Handle failure conditions first and return early.
- **Avoid `else`.** Refactor to early returns or ternaries.
- **Separate compound ifs.** Prefer nested `if` statements over `&&` chains.
- Always use curly brackets, even for a single statement.
- Ternary operators: keep each part on its own line unless the expression is very short.
- Add blank lines between statements so the code can breathe. The exception is a sequence of equivalent single-line operations.
- No extra empty lines between `{}` brackets.

```php
// Happy path last
if (! $user) {
    return null;
}

if (! $user->isActive()) {
    return null;
}

// Process active user...

// Short ternary
$name = $isFoo ? 'foo' : 'bar';

// Multi-line ternary
$result = $object instanceof Model ?
    $object->name :
    'A default value';

// Ternary instead of else
$condition
    ? $this->doSomething()
    : $this->doSomethingElse();

// Bad: compound condition with &&
if ($user->isActive() && $user->hasPermission('edit')) {
    $user->edit();
}

// Good: nested ifs
if ($user->isActive()) {
    if ($user->hasPermission('edit')) {
        $user->edit();
    }
}
```

### Comments

Be very critical about adding comments. They often become outdated and mislead over time. Code should be self-documenting through descriptive variable and function names. Adding a comment should never be the first tactic for making code readable.

*Instead of this:*

```php
// Get the failed checks for this site
$checks = $site->checks()->where('status', 'failed')->get();
```

*Do this:*

```php
$failedChecks = $site->checks()->where('status', 'failed')->get();
```

- Don't add comments that describe what the code does. Make the code describe itself.
- Short, readable code doesn't need comments explaining it.
- Use descriptive variable names instead of generic names plus a comment.
- Only add a comment to explain *why* something non-obvious is done, never *what* is being done.
- Refactor explanatory comments into descriptively named methods.
- Never add comments to tests. Test names should be descriptive enough.

---

## Laravel Conventions

### Configuration

- Config filenames: **kebab-case** (`media-library.php`, `pdf-generator.php`)
- Config keys: **snake_case** (`'chrome_path' => env('CHROME_PATH')`)
- Never use `env()` outside config files. Use the `config()` helper.
- Service-specific config goes in `config/services.php`, not a new file.

### Routing

- URLs: **kebab-case** (`/open-source`, `/front-end-developer`)
- Route names: **camelCase** (`->name('openSource')`)
- Route parameters: **camelCase** (`{newsItem}`, `{userId}`)
- HTTP verb first: `Route::get('open-source', [OpenSourceController::class, 'index'])`
- Use tuple notation `[Controller::class, 'method']`, not the string `'Controller@method'`
- Don't prefix URLs with `/`, except the root `/`

### API Routing

- Plural resource names: `/errors`, `/error-occurrences`
- Kebab-case resources
- Limit deep nesting. Prefer `/error-occurrences/1` over `/projects/1/errors/1/error-occurrences/1`
- Nest only when the context is necessary: `/errors/1/occurrences`

### Controllers

- **Plural** resource name plus a `Controller` suffix: `PostsController`
- Stick to the CRUD keywords: `index`, `create`, `store`, `show`, `edit`, `update`, `destroy`
- Extract new controllers for non-CRUD actions (for example `FavoritePostsController` with `store` and `destroy`)
- Use invokable controllers for single actions: `PerformCleanupController`

### Views & Blade

- View files: **camelCase** (`openSource.blade.php`)
- Indent with 4 spaces.
- No spaces after directives: `@if($condition)`
- Use `__()` for translations, not `@lang`

### Validation

- Always use array notation: `['required', 'email']`, never the pipe form `'required|email'`. Array notation is easier to combine with custom rule classes.

```php
public function rules() {
    return [
        'email' => ['required', 'email'],
    ];
}
```

- Custom rules use **snake_case**:

```php
Validator::extend('organisation_type', function ($attribute, $value) {
    return OrganisationType::isValid($value);
});
```

### Authorization

- Policies use **camelCase**: `Gate::define('editPost', ...)`
- Use CRUD words, but replace `show` with `view`

### Migrations

- Only write `up` methods in migrations. Don't write `down` methods.

### Artisan Commands

- Command names: **kebab-case** (`delete-old-records`)
- Always output feedback. At minimum, a `$this->comment('All ok!')` at the end.
- For batch processing, output progress per item and a summary at the end.
- Put the output *before* processing the item, which makes debugging a failure easier:

```php
$items->each(function(Item $item) {
    $this->info("Processing item id `{$item->id}`...");
    $this->processItem($item);
});

$this->comment("Processed {$items->count()} items.");
```

### Naming Classes

| Type | Convention | Example |
|------|-----------|---------|
| Controller | Plural + `Controller` | `PostsController` |
| Invokable Controller | Action + `Controller` | `PerformCleanupController` |
| Model | Singular | `Post` |
| Job | Action description | `CreateUser`, `SendEmailNotification` |
| Event | Tense indicates timing | `UserRegistering` / `UserRegistered` |
| Listener | Action + `Listener` | `SendInvitationMailListener` |
| Command | Action + `Command` | `PublishScheduledPostsCommand` |
| Mailable | Event/action + `Mail` | `AccountActivatedMail` |
| Resource | Plural + `Resource` | `UsersResource` |
| Enum | Descriptive, no prefix | `OrderStatus`, `BookingType` |

---

## Common Mistakes to Avoid

- Using `env()` outside config files
- Using pipe notation for validation rules (`'required|email'`)
- Using `$fillable` instead of `$guarded = []` in package models
- Adding spaces after Blade directives (`@if ($condition)`)
- Putting extra empty lines inside `{}` brackets
- Using `final` on classes (Spatie doesn't by default)
- Docblocks on fully type-hinted methods without descriptions
- Fully qualified class names inside docblocks instead of an imported short name
- Inline fully qualified class names in code instead of a `use` statement
- Single-letter variable names such as `$e`
- String controller references (`'Controller@method'`) instead of tuple notation
- Using `else` where early returns work
- Deep API route nesting when a flat route suffices
- Creating new config files for service credentials (use `services.php`)
- Forgetting the `void` return type on methods that return nothing
- Writing `down` methods in migrations
- Introducing a `private const` instead of inlining it or using a private property

---

## Detailed references

Load these as needed:

- **Package architecture**: `references/package-architecture.md`. Structure, service providers, contracts, model and config patterns.
- **Testing with Pest**: `references/testing-pest.md`. Testbench setup, test style, what to test, helpers, composer stack.
- **Git workflow**: `references/version-control.md`. Branch naming, PR workflow, commit conventions.
- **Laravel and PHP style, in depth**: `references/laravel-php.md`.
- **JavaScript style**: `references/javascript.md`.
- **New project setup**: `references/new-project-setup.md`.
