---
name: explain-changes
description: Explain how the code on the current branch works, as an HTML page opened in the browser, written so someone can understand the change and review it. Use when asked to explain a branch, explain these changes, walk me through this code, help me review this, or explain what an agent just did. Most useful when someone other than the reader wrote the code and it has to be understood before merging.
---

# Explain changes

Produce a detailed written explanation of everything on this branch, as an HTML page, and open it.

The page has one job: let the reader understand the change well enough to review it. Those two things are not separate. Nobody can review code they do not understand, and explaining code properly is what surfaces the problems in it.

So write the explanation, and when writing it reveals something wrong, say so. You cannot write "this is correct because the guard above catches the null case" without going to look at whether that guard exists. When it does not, you have found a real bug, and it came from understanding rather than from a checklist. If you cannot explain why a line is correct, that inability is itself the finding.

Do not fix anything and do not approve anything. Do not run a systematic bug hunt over the diff either; that is a separate task with its own checklists. Explain the change, and report what the explaining turned up.

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

## How to write it

Explain the change the way you would to a colleague sitting next to you, in whatever order actually explains it. Follow the mechanism, not the file list. There is no required set of sections: decide what this particular change needs, and write that.

- **Show the code.** Real excerpts from the change and from the existing code around it, with the file path above them. Never invent an illustrative example. If you cannot show the line, do not make the claim.
- **Explain, do not summarise.** "Adds caching to the report query" is a commit message. Explain what is cached, keyed on what, invalidated when, and what happens on a miss.
- **Explain the choices.** Where one approach was taken over an obvious alternative, name the alternative. A diff shows what was done and never what was rejected, and that fork in the road is where agent-written code usually goes wrong: a new helper where one already existed, an observer where the codebase uses events, a pattern the rest of the project deliberately avoids.
- **Point at what deserves scrutiny.** Authorization, money, data loss, migrations, multi-tenancy, anything that fails quietly. Also what did not change but arguably should: call sites not updated, tests not added, documentation now stale.
- **Say when you do not know.** Write "unclear why this was needed" rather than inventing a rationale. An honest gap tells the reader where to look. A confident guess sends them past the problem.
- **Never explain around a problem.** The temptation, when code does not make sense, is to write a plausible paragraph that makes it sound deliberate. That is the single worst thing this page can do, because it launders a bug into an explanation and the reader stops looking. Say it does not make sense.
- Anything the explanation turned up goes inline, next to the code it concerns, not collected into a list at the bottom. The reader needs the surrounding code in view to judge it.
- **Aim for readable in ten minutes** for a change of a few hundred lines. For much larger changes, explain the mechanism thoroughly once and then summarise the repetitive applications of it.
- Avoid em dashes and en dashes as punctuation. Periods, commas, and parentheses read better.

## Link to the code

If the branch is pushed, make every file path and code block heading a link to the real thing, so the reader can jump from the explanation to the code and to the discussion around it.

```bash
git remote get-url origin                       # the repo
git rev-parse HEAD                              # pin links to this commit
git branch --remote --contains HEAD             # is it actually pushed?
gh pr view --json url,number 2>/dev/null        # is there a PR?
```

Build permalinks against the commit sha rather than the branch name, so they keep pointing at the code being explained after the branch moves on:

```
https://github.com/<owner>/<repo>/blob/<sha>/<path>#L42-L60
```

Link the PR itself in the header when there is one. If the branch is not pushed, say so once in the header and use plain paths. Never emit a link you have not confirmed resolves to a pushed commit; a 404 in a review page is worse than no link.

## Output

Write the page outside the repository:

```
~/.claude/explanations/<repo>-<branch>-<YYYYMMDD-HHMMSS>.html
```

**Never write it inside the working tree**, not even untracked. It would show up in `git status` while the reader is trying to read a clean diff, and any tooling that commits pending work before merging would carry it into the default branch.

Build the page from `references/page-template.html`, which carries the styling and loads highlight.js.

Write code blocks as plain code with the language set on the element, `<code class="language-php">`. Do not mark up tokens by hand; highlight.js colours them. Use `language-diff` with real `+` and `-` lines when a before and after belong side by side.

Then open it and print the path:

```bash
mkdir -p ~/.claude/explanations
open "$OUTPUT"
```
