---
name: audit-architecture
description: Run an application-wide, read-only architecture audit that finds materially useful simplifications in a codebase's data structures, state representation, control flow, algorithms, and ownership boundaries. Fans out bounded read-only agents per subsystem, verifies every finding against the repo, and ranks results P0 to P3. Use when asked to audit architecture, review how state or data is modelled, find structural simplifications, or invoke /audit-architecture. Not a style or line-level quality pass (use /simplify for that) and not a bug hunt (use /code-review).
disable-model-invocation: true
---

# DSA Codebase Audit

Audit this entire codebase for materially useful simplifications in its data structures, state representation, control flow, algorithms, and ownership.

This is an audit-only exercise. Do not edit files, run tests, implement recommendations, commit, or push. Read-only inspection commands are allowed.

You are the coordinator. Continue until the complete codebase has been reviewed and the final audit is validated.

## 1. Establish the coverage contract

Inspect the repository and inventory every identifiable subsystem.

Give each subsystem:
- a stable ID and descriptive name;
- an exact ownership boundary;
- its key implementation files;
- relevant public interfaces, major call sites, and tests;
- a status: queued, in review, recommend, or skip.

Include frontend, backend, shared infrastructure, platform bridges, generated-contract ownership, and test/tooling infrastructure where materially relevant.

Create one canonical scratchpad or report containing:
- the subsystem inventory;
- confirmed opportunities;
- explicit skip decisions;
- cross-cutting patterns;
- duplicates and superseded findings;
- final priorities and dependencies;
- an audit log.

Treat this inventory as the coverage contract. Do not assume broad catch-all rows prove coverage.

## 2. Run bounded subsystem reviews

Use fresh, read-only agents where available. Give every worker one distinct subsystem with an exact, non-overlapping ownership boundary.

Keep concurrency bounded to the number of lanes you can actively coordinate. Use one consolidated wait mechanism, do not interrupt productive workers merely because they are slow, and close completed workers after harvesting their results.

Each worker receives this brief:

> Review the assigned subsystem for at most two materially useful simplifications in its data structures, state representation, or organizing model.
>
> Inspect its implementation, public interfaces, major call sites, and existing tests. Stay within the assigned ownership boundary. You may identify cross-subsystem concerns, but do not expand the scope to solve them.
>
> Look for:
> - scattered booleans or nullable fields that permit invalid combinations and should become a state machine or discriminated union;
> - repeated assumptions about object shape that need a shared typed model;
> - duplicated branching that a small map, registry, reducer, or command model would remove;
> - unclear state or behavior ownership that a small module boundary would clarify;
> - repeated scans, transformations, or lookups where a more appropriate collection or index would materially simplify behavior;
> - lifecycle, concurrency, or async states whose representation permits stale or contradictory state.
>
> Do not force an abstraction. Prefer boring local code when it is already clear.
>
> Do not recommend changes solely for stylistic consistency, hypothetical extensibility, minor line-count reduction, or moving existing branching behind a new type.
>
> Return at most two opportunities. If nothing clearly meets the threshold, return `skip`.
>
> For every recommendation, provide:
>
> 1. Verdict: recommend or skip.
> 2. Evidence with exact file and line references.
> 3. Current complexity or invalid states.
> 4. Proposed representation and why it is simpler.
> 5. Smallest credible implementation scope, including affected files and interfaces.
> 6. Regression risks and migration concerns.
> 7. Existing and additional validation required.
> 8. Confidence: high, medium, or low.

## 3. Validate and synthesize

The coordinator must independently verify every finding against the current repository before accepting it.

Reject, narrow, or demote recommendations that are vague, duplicate another finding, misunderstand intentional semantics, or merely relocate complexity.

Record skips as completed coverage. Deduplicate overlapping findings and assign each accepted recommendation to one authoritative subsystem.

Continue opening bounded review batches until every inventory row is complete.

## 4. Audit the audit

Before finishing, run fresh independent passes for:

- repository coverage and missing subsystem boundaries;
- duplication and ownership overlap;
- materiality and over-abstraction;
- schema completeness;
- dependency-aware priority ranking.

If the coverage pass finds a real omission, add an explicit subsystem row and audit it. Do not hide it by broadening a previously completed boundary.

Rank the final recommendations by concrete impact, confidence, implementation effort, blast radius, and prerequisites. Identify the best first implementation slices.

The audit is complete only when:

- every identifiable subsystem has been reviewed;
- every subsystem has a recommendation or explicit skip;
- every finding has complete evidence, scope, risk, and validation fields;
- duplicates and weak abstractions have been removed;
- priorities and dependencies are internally consistent;
- the repository remains unchanged.

## Final output: priority synthesis

Close with a priority synthesis: the total number of unique recommendations, a guarantee that every promoted ID appears exactly once, and the IDs grouped into tiers.

- **P0**: correctness, security, or data-loss risk. Reachable wrong-record, lost-update, authorization, durable-state, or permanently incomplete-operation risks.
- **P1**: concrete correctness or high-leverage contract work. Concrete boundary failures and high-leverage ownership fixes with less immediate damage or greater migration cost.
- **P2**: material invariant improvements whose observed impact is narrower or whose migration is sensitive.
- **P3**: lower-impact telemetry, diagnostics, or maintainability improvements. Expect this tier to be small after materiality review.
