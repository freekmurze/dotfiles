# Review lanes

The shared review lanes used by `review-code` (working tree) and `review-pr` (a GitHub PR).

Run every applicable lane concurrently, then aggregate.

## Lane 1: correctness (always)

Find bugs in the target. This is the only lane that looks for correctness. The convention lanes below check style and project conventions, and will happily pass code that is flawlessly conventional and functionally wrong.

## Lane 2: PHP simplification

For changed PHP files, simplify and refine the PHP and Laravel code.

## Lane 3: Spatie conventions

For changed PHP files, check compliance with the `spatie-guidelines` skill.

## Lane 4: Laravel best practices

For changed PHP files, check routing, database performance, and architecture against the `laravel-best-practices` skill.

## Lane 5: React

Only when JS or TS files changed. Check against the `react-best-practices` skill.

## Running the lanes

### Claude Code

Run all lanes in parallel, in a single message.

| Lane | How |
|---|---|
| 1 correctness | `/code-review` on the target |
| 2 PHP simplification | the `laravel-simplifier:laravel-simplifier` agent |
| 3 Spatie conventions | a `general-purpose` agent invoking `spatie-guidelines` |
| 4 Laravel practices | a `general-purpose` agent invoking `laravel-best-practices` |
| 5 React | a `general-purpose` agent invoking `react-best-practices` |

### Codex

| Lane | How |
|---|---|
| 1 correctness | the `review-agent` skill, delegating the target to it |
| 2 PHP simplification | the simplifier is a Claude plugin agent, but its definition is a plain file. Read `~/.claude/plugins/cache/laravel/laravel-simplifier/*/agents/laravel-simplifier.md` and apply it yourself. If that path does not exist, say the lane was skipped. |
| 3 Spatie conventions | read `spatie-guidelines` and apply it to the changed files yourself |
| 4 Laravel practices | **only if `laravel-best-practices` is readable from this repo.** It is normally installed by Laravel Boost into `<repo>/.claude/skills/`, which Codex does not load. Read that path directly if it exists, otherwise say the lane was skipped. |
| 5 React | read `react-best-practices` and apply it to the changed files yourself |

Never report a review as complete without naming the lanes that were skipped.

### Any other harness

Lane 1 is a defect-first review. Lanes 2 to 5 are each "read the named skill, apply it to the changed files, report deviations". Use whatever delegation the harness offers, or do them inline.

## Brief for every lane

- Review only the files in scope. Never touch unrelated WIP files in the working tree.
- Report findings as a punch list with `file:line` references.
- Do not fix anything during the lane. Report only.

## Aggregating the results

Deduplicate across lanes. When two lanes flag the same line, keep the correctness finding and drop the convention one.

Correctness findings take priority. Never apply a convention fix that changes behaviour flagged as a bug. Fix the bug first.

Group the summary as **Correctness** and **Conventions**. No "Test plan" section.
