---
name: review-pr
description: Review and merge GitHub pull requests for Spatie packages. Use when asked to review a PR, review a pull request, merge a PR, or when given a GitHub PR URL to review. Also triggers on 'review this PR,' 'check this pull request,' 'merge this,' or '/review-pr'. Uses gh CLI for all GitHub operations.
---

# Review PR

Review a GitHub pull request, verify CI status, merge it, and thank the author. Only tag a release when explicitly asked.

## Workflow

### 1. Fetch PR details

Use `gh` CLI to get the PR diff, description, and CI status:

```bash
gh pr view <number> --json title,body,additions,deletions,files,reviews,statusCheckRollup,headRefName,baseRefName
gh pr diff <number>
```

### 2. Review the changes

Run the shared review lanes against this PR, in **report-only mode**. Read the lane definitions from:

```
~/.claude/skills/review-code/references/lanes.md
```

Run every applicable lane against the PR, following the "Running the lanes" section for whichever harness you are. Give each lane the diff from step 1 as its target.

**Report only. Never apply fixes to someone else's PR**, even for findings that would be auto-applied in `/review-code`. The output of this step is a verdict, not an edit.

On top of the lane findings, judge:

- Whether the change matches the PR description
- Test coverage for new functionality
- Whether it fits the existing patterns of this package

If there are issues, post a review comment via `gh pr review <number> --request-changes --body "..."` and stop.

### 3. Check CI status

All CI checks must be green before merging. Verify via the `statusCheckRollup` field from step 1. If CI is failing or pending, inform the user and stop.

### 4. Merge

If the review looks good and CI is green, merge the PR:

```bash
gh pr merge <number> --squash --delete-branch
```

### 5. Thank the author

After merging, thank the author of the work in the conversation this change came from. Post the thanks on whichever of these the change originated from, in this order of preference:

1. The issue the PR closes, if it references one.
2. The discussion the PR came out of, if there is one.
3. The PR itself.

Find the author of *that* conversation, not necessarily the PR author. They are often different people: someone reports an issue, someone else fixes it. Thank both when they differ.

```bash
gh issue comment <number> --body "..."
gh pr comment <number> --body "..."
```

Keep it short, warm, and specific about what they contributed. Never mention Claude Code.

### 6. Release, only when I ask

**Never tag a release unless I explicitly ask for one in this conversation.** Merging is not a request to release. Finishing a review is not a request to release. Say the PR is merged and stop there.

When I do ask, determine the version bump by checking the latest tag:

```bash
gh release list --limit 1
```

#### Versioning rules

- **NEVER tag a major version bump without explicit user approval.** Always ask first.
- Dropping support for a PHP or Laravel version is NOT a breaking change, use a minor or patch bump.
- New features: minor version bump (e.g. 1.2.0 to 1.3.0)
- Bug fixes, typos, dependency updates: patch version bump (e.g. 1.2.0 to 1.2.1)
- When in doubt about the version bump level, ask the user.

#### Creating the release

```bash
gh release create <tag> --title "<tag>" --generate-notes
```

Use `--generate-notes` to auto-generate release notes from the merged PR.
