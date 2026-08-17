---
name: explain-changes
description: Explain how the code on the current branch works, as a self-contained HTML page opened in the browser, written to help someone review it. Use when asked to explain a branch, explain these changes, walk me through this code, help me review this, or explain what an agent just did. Most useful when someone other than the reader wrote the code and it has to be understood before merging.
---

# Explain changes

Produce a detailed written explanation of everything on this branch, as an HTML page, and open it.

The goal is comprehension, not judgment. Do not fix anything, do not approve anything, and do not run a systematic bug hunt over the diff. That is a separate task with its own checklists.

**But say what the explaining reveals.** Writing an explanation is a stronger comprehension test than reading a diff. You can skim a diff and nod. You cannot write "this is correct because the guard above catches the null case" without going to look at whether that guard exists. When it does not, you have found a real bug, and it came from understanding rather than from a checklist.

So when explaining a piece of code shows that it is wrong, say so plainly and prominently, then carry on explaining. If you cannot explain why a line is correct, that inability is itself the finding. Write it down rather than smoothing over it.

## Scope

Everything on this branch that is not in the default branch, **committed and uncommitted**:

```bash
BASE=$(git merge-base main HEAD)   # or master, whichever the repo uses
git diff "$BASE"..HEAD             # committed
git diff HEAD                      # uncommitted
git status --short                 # untracked files
```

Include uncommitted work. An agent often leaves changes uncommitted, and merge tooling frequently commits pending work on the way through, so uncommitted changes ship. A page that silently skipped them would describe less than what is about to be merged.

On the default branch with no feature branch, fall back to the working tree alone and say so on the page.

Only explain something else if the user asks for it.

## Read beyond the diff

A diff shows `$query->where('workspace_id', $id)` but not whether a global scope already applied that. Read the surrounding file, the class it extends, and the call sites, enough to explain how the code actually behaves. An explanation assembled from the diff alone is something `git diff` already provides.

## What the page is for

The reader is about to review this change. They need to understand the mechanism, know where the risk is, and see the choices that were made. A structure that serves that:

| Section | What it is for |
|---------|----------------|
| What changed and why | The whole change in a few sentences, plus the stats. Someone who reads only this should know what they are dealing with |
| How it works | The mechanism, following the execution path rather than the file list. Entry point, what happens, where it lands |
| Decisions made | Where one approach was chosen over another, and what the alternative was. A diff shows what was done and never what was rejected |
| What to look at hardest | Authorization, money, data loss, migrations, multi-tenancy, anything that fails quietly |
| What did not change but arguably should | Call sites not updated, tests not added, documentation now stale |
| File index | Orientation for a large change |

Anything the explanation turned up goes inline, at the point in the walkthrough where it became apparent, not collected into a list at the bottom. The reader needs the surrounding code in view to judge it.

**This structure is a suggestion, not a template.** Reorganise it, merge sections, or invent different ones whenever that explains the change better. Judge any deviation against the purposes above, not against the layout.

The "Decisions made" section is the one most worth keeping. When an agent wrote the code, the failure mode is rarely a syntax error. It is a reasonable-looking choice that is wrong for this codebase: a new helper where one already existed, an observer where the codebase uses events, a pattern the rest of the project deliberately avoids. The fork in the road is invisible in a diff.

## Rules for the writing

- **Show real code.** Excerpts from the change and from the existing code around it. Never invent an illustrative example. If you cannot show the line, do not make the claim.
- **Say when you do not know.** Write "unclear why this was needed" rather than inventing a rationale. An honest gap tells the reader where to look. A confident guess sends them past the problem.
- **Never explain around a problem.** The temptation, when a piece of code does not make sense, is to write a plausible paragraph that makes it sound deliberate. That is the single worst thing this page can do, because it launders a bug into an explanation and the reader stops looking. Say it does not make sense.
- **Explain, do not summarise.** "Adds caching to the report query" is a commit message. Explain what is cached, keyed on what, invalidated when, and what happens on a miss.
- **Explain the choice, not just the code.** Where one approach was taken over an obvious alternative, name the alternative. A reader can evaluate a decision far faster than they can reconstruct it.
- **Aim for readable in ten minutes** for a change of a few hundred lines. For much larger changes, explain the mechanism thoroughly once and then summarise the repetitive applications of it rather than marching through every instance.
- Avoid em dashes and en dashes as punctuation. Periods, commas, and parentheses read better.

## Output

Write the page outside the repository:

```
~/.claude/explanations/<repo>-<branch>-<YYYYMMDD-HHMMSS>.html
```

**Never write it inside the working tree**, not even untracked. It would show up in `git status` while the reader is trying to read a clean diff, and any tooling that commits pending work before merging would carry it into the default branch.

Build the page from `references/page-template.html`, which carries the styling. Keep it self-contained: no CDN, no remote fonts, no external scripts. It is opened over `file://` and has to work offline.

Then open it and print the path:

```bash
mkdir -p ~/.claude/explanations
open "$OUTPUT"
```
