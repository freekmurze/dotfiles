---
name: review-code
description: Review all changed code against project conventions
disable-model-invocation: true
---

Review the code I just changed, in the working tree. Do NOT touch unrelated WIP files.

## Scope

Everything changed in the working tree (staged and unstaged), relative to HEAD.

## Run the lanes

Run every applicable lane from `references/lanes.md`, in parallel, in a single message.

## Apply the findings

Aggregate and deduplicate as described in `references/lanes.md`, then act on each finding:

- Clear correction with a concrete fix, apply it directly with Edit.
- Judgment call or false positive, note it in the summary and skip.

End with a short bulleted summary grouped as **Correctness** and **Conventions**, covering what was changed (or confirming the code was already clean). No "Test plan" section.
