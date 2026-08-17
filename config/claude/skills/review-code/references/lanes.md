# Review lanes

The shared review lanes used by `/review-code` (working tree) and `/review-pr` (a GitHub PR).

Run every applicable lane **in parallel**, in a single message.

## Lane 1: correctness (always)

Run `/code-review` on the target. This is the only lane that looks for bugs. The convention lanes below check style and project conventions, and will happily pass code that is flawlessly conventional and functionally wrong.

## Lane 2: PHP simplification

For changed PHP files, use the `laravel-simplifier:laravel-simplifier` agent to simplify and refine PHP and Laravel code.

## Lane 3: Spatie conventions

For changed PHP files, use a `general-purpose` agent invoking the `spatie-guidelines` skill to check compliance with Spatie's PHP and Laravel guidelines.

## Lane 4: Laravel best practices

For changed PHP files, use a `general-purpose` agent invoking the `laravel-best-practices` skill to check routing, database performance, and architecture.

## Lane 5: React

Only when JS or TS files changed. Use a `general-purpose` agent invoking the `react-best-practices` skill.

## Brief for every lane agent

- Review only the files in scope. Never touch unrelated WIP files in the working tree.
- Report findings as a punch list with `file:line` references.
- Do not fix anything. Report only.

## Aggregating the results

Deduplicate across lanes. When two lanes flag the same line, keep the correctness finding and drop the convention one.

Correctness findings take priority. Never apply a convention fix that changes behaviour flagged as a bug. Fix the bug first.

Group the summary as **Correctness** and **Conventions**. No "Test plan" section.
