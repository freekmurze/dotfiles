---
name: designing-spatie-rust-ports
description: "Designs Spatie Rust APIs when porting PHP packages. Use when creating or reviewing a Spatie Rust port, mapping fluent APIs, documenting PHP behavior differences, or adapting the Rust package skeleton. Complements rust-best-practices; not a general Rust style guide."
---

# Designing Spatie Rust ports

Preserve the PHP package's vocabulary, defaults and user journey. Change the language mechanics where Rust requires it, and record observable behavior differences rather than hiding them behind familiar method names.

Use `rust-best-practices` for ownership, error implementation, generics and ordinary Rust quality rules. This skill decides how those rules meet Spatie's PHP APIs. Follow an existing crate's documented contracts when maintaining it; these defaults are not permission to rename a released API or raise its minimum Rust version.

## Establish the PHP contract first

Read the PHP implementation, tests and documentation at a specific tag or commit. Read existing Rust callers and tests before changing an established port. Inventory public behavior, not just public method names: defaults, validation, output, ordering, mutation, callbacks, cleanup and failure policy.

Keep a parity record, meaning a mapping of PHP behavior to its Rust equivalent and verification. For each feature, record the PHP reference revision, Rust API, status, deliberate difference or omission, and a test or fixture that establishes the claim. Distinguish implemented-and-executed, compiled-only, unverified and intentionally omitted behavior. A passing adapter test does not establish compatibility with a real external runtime.

Update an existing parity document when the port has one. Otherwise use a concise maintainer section or focused usage notes; add a separate document only when the mapping needs it or the task requests it. Do not add mandatory parity files to the generic skeleton. Keep differences that affect users beside the relevant usage example, even when a longer mapping exists elsewhere.

## Set up every new package

1. Freek creates the public `spatie/<php-package>-rust` repository from `spatie/package-skeleton-rust`; Amp's GitHub integration cannot create it. Give him `https://github.com/new?template_owner=spatie&template_name=package-skeleton-rust&owner=spatie&name=<php-package>-rust&visibility=public`, replacing `<php-package>` with the PHP repository name. Keep the skeleton's packaging and documentation machinery; replace its sample API and tests.
2. Disable Wiki and Projects. Copy the PHP repository's description and website exactly, and copy its topics plus `rust`. Under "Include in the home page", enable only Releases; turn Deployments and Packages off. These home-page checkboxes are UI-only: use a signed-in browser on Freek's Mac, not an API approximation.
3. Set the README H1 from the PHP README title using the rule under "Make the docs familiar and executable".
4. Get CI green, then import the docs on spatie.be. With explicit release approval, publish version `1.0.0` to crates.io and add `github:spatie:crates` as an owner. This checklist is not authorization to publish or change shared settings.

## Keep names recognizable

- Cargo packages are `spatie-<name>`, such as `spatie-calendar-links`. Rust imports use underscores, such as `spatie_calendar_links`. Repository names are separate, commonly `<name>-rust`; do not derive the crate name blindly from the repository name.
- Keep domain types such as `Period`, `Link`, `SslCertificate` and `TemporaryDirectory`. Convert method names to snake_case: `usingPort` becomes `using_port`.
- Preserve recognizable named constructors such as `make`, `create` and `download` when they express the PHP workflow. Use `new` for ordinary construction; do not force every PHP factory into `new` or add aliases for every spelling.
- Follow the port's naming convention for accessors. `SslCertificate::issuer()` drops PHP's `get`, while `DbDumper::get_db_name()` retains it. Both exist deliberately; do not run a blanket `get_` or `set_` removal across the API.
- Replace PHP's optional-argument overloads with explicit Rust operations when that reads better: `path()` / `path_for(child)` and `optimize(path)` / `optimize_to(input, output)`. Use `Option<T>` when one operation naturally has an optional value. Record the mapping.

## Choose fluent ownership by what the object does

For a value configured before use, default to setters that take `self` and return `Self`. A setter that validates its input may return `Result<Self, Error>`. The caller chains the methods or assigns the returned object. Mark consuming builders or methods `#[must_use]` so an ignored result is visible.

Keep the PHP sequence recognizable. Do not insert `.build()` when `create()`, `save()` or another existing terminal operation already does the work. Add a separate build step only when it creates a distinct validated object callers need to retain.

Use `&mut self` and, when useful, `&mut Self` for objects callers continue editing. Schema.org builders are an important exception to consuming chains: they preserve shared entity identity, so later edits can affect graph references and cloning shares an entity. Preserve those semantics and document them; changing ownership must not silently turn a shared graph into disconnected copies.

Value APIs stay value APIs. A period comparison should not acquire mutable builder state merely because other Spatie packages are fluent. Use traits or closures at actual extension points, such as optimizers or generators; do not reproduce PHP inheritance, runtime macros or reflection with a Rust compatibility layer.

## Translate exceptions and configuration into explicit Rust contracts

Map recoverable PHP exceptions to `Result<T, Error>` and ordinary absence to `Option<T>`. Retain distinctions callers need, such as invalid input, process failure and timeout. Do not expose an untyped string or a boolean in place of a failure the caller must handle.

Do not add `Result` to infallible operations for symmetry. For configuration that can be fixed by a later setter, validate the combined settings at preparation or execution. Validate independently invalid values at the setter when appropriate. Document where failure occurs so PHP users know where to put `?`.

Preserve the PHP failure policy as well as its exception cases. For example, image optimization can log an optimizer failure and continue by default; `throws(true)` changes that to returning an error, and a Rust callback returns `Result` to choose. Do not reinterpret such a policy as catching Rust panics or turning every failure into `Ok`.

Keep defaults aligned unless a Rust/runtime or safety difference is explicit. Use typed configuration for structured options instead of reproducing PHP config arrays: enums for finite choices, `Duration` for durations where the API adopts it, path types for paths, and types such as `u16` for ports. Preserve familiar units in existing APIs and make conversions explicit; a PHP seconds setter must not quietly become milliseconds.

Keep configuration owned by the caller or an explicit runtime context. Pass dependencies or expose traits/closures instead of introducing hidden global services. A standalone package does not need service registration or a framework emulation layer. Preserve an established async boundary, but do not require an async runtime merely to imitate fluent PHP syntax.

## Preserve useful behavior, not obsolete plumbing

Port the feature's purpose rather than its PHP implementation dependencies. Browsershot controls Chrome directly, so Node/Puppeteer executable and module-path configuration does not belong in its Rust API. Omit such options and explain the requirement change; do not leave successful no-op setters that suggest they work.

This is not a blanket ban on Node or external tools. Image optimization can still invoke SVGO, whose installation requires Node. Document requirements per feature and backend.

Use native Rust process arguments rather than PHP shell fragments, and explain split arguments, unsupported pipes and redirects. Use ownership-based cleanup for temporary resources; document whether cleanup is automatic or opt-in and how callers obtain cleanup errors. Do not promise that `Drop` reports errors or that a failed operation leaves output untouched unless the implementation guarantees it.

Do not silently copy an unsafe behavior or silently change a domain rule in the name of safety. Document concrete safety improvements and their compatibility effects, and keep the original domain behavior unless the task calls for a change. For example, a certificate's metadata predicate is not automatically a cryptographic trust check.

## Make the docs familiar and executable

Mirror the PHP guide's headings, order, examples, variable names, values and wording where the behavior still matches. If the PHP package has only a README, use its usage sections as the guide source. Change prose that would become false, including installation, exceptions, runtime requirements and safety guarantees.

Start from the PHP package's descriptive README H1, then adapt it for Rust readers rather than copying it literally. Remove package-name prefixes and replace PHP with Rust where it describes the user's language. For example, "spatie/ray - Send debug output from PHP to Ray" becomes "Send debug output from Rust to Ray". Keep PHP references that identify the upstream package or compatibility differences. Check the README H1, docs introduction headings, crate description and docs slogan for the same mistakes. If the PHP title would make a false runtime claim, use the PHP repository's GitHub description and adapt it by the same rule. Browsershot uses "Convert HTML to an image, PDF or string", not "spatie-browsershot" or a Puppeteer claim. Supply the resulting descriptive title to the skeleton's configure step.

Use the Rust skeleton's README structure and badges, with Cargo installation and testing commands. Keep the original package and author credits. README ad URLs use the repository key, such as `period-rust`, not the crate key `spatie-period`; check that the image actually exists. Keep `docs/_index.md` and page frontmatter in the skeleton's Spatie-site format. Website registration and hosting are separate from writing Markdown.

Keep the visible example focused on the PHP-sized workflow. Hide imports, error-return plumbing, isolated fixtures, cleanup and assertions on rustdoc `# ` lines. Setup must make the shown operation real; never hide a replacement implementation or the operation being demonstrated.

For example, the temporary-directory workflow can remain this short while still being an executable doctest:

```rust
# fn main() -> Result<(), spatie_temporary_directory::Error> {
use spatie_temporary_directory::TemporaryDirectory;

let directory = TemporaryDirectory::new()
    .delete_when_destroyed(true)
    .create()?;
# assert!(directory.path().is_dir());
# Ok(())
# }
```

Run self-contained examples and assert claimed outputs using deterministic fixtures. A certificate example claiming `spatie.be` needs a certificate fixture with that subject; an unrelated PEM file only proves parsing. Label synthetic fixtures honestly.

Use `rust,no_run` only for examples that genuinely require external credentials, services or tools unavailable in the ordinary doctest environment. They compile but do not execute; cover their behavior in separate integration checks and state what was not run. Do not use `ignore` or `no_run` to conceal broken self-contained examples.

Retain README rustdoc inclusion and recursive guide doctest discovery when supplied by the skeleton. Include the build script, guides and referenced fixtures in the Cargo package. Run doctests separately because `--all-targets` does not include them. Check rendered docs too: rustdoc hides `# ` lines, but a website renderer must support that convention independently. Hidden lines in a Markdown file do not prove they are hidden on the website.

## Keep the toolchain and CI bar consistent

For new ports, use Rust edition `2024` and `rust-version = "1.89"`. The minimum supported Rust version, or MSRV, is tested with Rust `1.89.0`; it is not inferred from success on stable. Keep the lockfile committed and check dependency changes against the MSRV.

Start with the skeleton's manifest lints: forbid unsafe code, warn for missing public docs and missing public `Debug`. Keep stricter settings in existing crates. CI denies Rust and rustdoc warnings. Do not relax existing lints to fit this draft or expose secrets just to derive `Debug`; use a redacted implementation or a narrowly documented exception.

Use the skeleton's stable Linux/macOS/Windows checks, a separate MSRV job, security audit and package verification. The core commands are:

```sh
cargo fmt --all -- --check
cargo clippy --locked --all-targets --all-features -- -D warnings
cargo test --locked --all-targets --all-features
cargo test --locked --doc --all-features
cargo doc --locked --no-deps --all-features
cargo package --locked
cargo +1.89.0 test --locked --all-targets --all-features
cargo +1.89.0 test --locked --doc --all-features
cargo audit
```

Set `RUSTFLAGS=-D warnings` and `RUSTDOCFLAGS=-D warnings` for the CI jobs. Add `--workspace` for a Cargo workspace; check separate manifests separately. Use the crate's feature matrix when feature combinations are mutually exclusive. The skeleton uses `cargo package --locked --allow-dirty` in CI to verify its generated/configured checkout; this is packaging verification, not publication.

Add native-tool, service or generator checks when those are part of the package contract. A fake executable can prove argument handling, but not a successful real dump and restore. Generated APIs need a reproducible generator and a no-diff generation check. Run tests against disposable inputs and record external requirements and unverified platforms.

## Keep public material limited to public facts

Derive public examples and source links from public repositories. Keep private repository names, unreleased runtime plans and private integration details out of public skills, docs, READMEs, commit messages and package metadata. Describe generic application contracts without revealing private projects. Review provenance before moving a draft to a shared or public location.

## Review a port against its PHP story

Before calling the port ready, walk the PHP examples against the Rust API, check the parity record, and run the applicable CI checks. Confirm defaults and errors, not only happy-path syntax. Confirm that shared objects still share, consuming setters are retained, external requirements are visible, and every omission has a reason. Report separately what compiled, what executed locally, what passed on CI and what remains unverified.
