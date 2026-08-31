# 17: Docs and GitHub issue templates

**What to build:** Bring the committed docs in line with what was built and add the issue templates the Help links point at. `DEPENDENCIES.md`: replace the "Future Dependencies" shadcn/ui and Wails entries with the real stack (Wails v3, React, TypeScript, Tailwind v4, Radix primitives, `lumberjack`) and record why shadcn was not used. Root `CLAUDE.md`: update the frontend section to match. `README.md`: update the technology-stack table and the development-status checkboxes. Add `.github/ISSUE_TEMPLATE/bug_report.md` and `feature_request.md`. Keep to the project rule against naming or comparing to the reference implementation in committed docs.

**Blocked by:** 11, 09 (enough of the build to describe accurately); run last

**Status:** done

- [x] `DEPENDENCIES.md` reflects the chosen frontend stack and drops the stale "future" entries; one line on why not shadcn
- [x] Root `CLAUDE.md` frontend section matches the built stack
- [x] `README.md` stack table and status list updated
- [x] `.github/ISSUE_TEMPLATE/bug_report.md` and `feature_request.md` exist and match the Help-section links
- [x] No committed doc names or compares against the reference implementation
- [x] `CONTEXT.md` glossary re-read for accuracy against the final code; fix any drift

**Notes:** `CONTEXT.md` had no drift; whoever carried tickets 04/11-16 already kept it current
(`Close Action`, `Onboarding`, and the removal of `Mode Override` were all already reflected).
`README.md` and `DEPENDENCIES.md` were still describing the original plan (shadcn/ui, "to be
added" frontend, a flat `Config`) and needed real edits. Root `CLAUDE.md` is a broad
pre-implementation planning doc; only the shadcn/ui and frontend-stack lines were corrected,
not the whole speculative file (e.g. the Huma API layer it describes was never built and is out
of scope here). `.github/ISSUE_TEMPLATE/` filenames match the `?template=` query params already
hardcoded in `frontend/src/lib/links.ts`. Per the user's instruction, these templates were added
to the current repo (the `Its-Haze/league-rpc-gui` release-testing remote, per
`internal/updates/config.go`'s `RepoSlug`) to see how they read; the Help-section links
already point at the eventual `its-haze/league-rpc` home and are untouched. The `spec.md`
Progress table's stale ticket-10 checkbox (already done, just never checked off after the later
rework commits) was corrected alongside this ticket's own.
