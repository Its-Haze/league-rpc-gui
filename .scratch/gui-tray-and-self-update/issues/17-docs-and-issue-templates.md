# 17: Docs and GitHub issue templates

**What to build:** Bring the committed docs in line with what was built and add the issue templates the Help links point at. `DEPENDENCIES.md`: replace the "Future Dependencies" shadcn/ui and Wails entries with the real stack (Wails v3, React, TypeScript, Tailwind v4, Radix primitives, `lumberjack`) and record why shadcn was not used. Root `CLAUDE.md`: update the frontend section to match. `README.md`: update the technology-stack table and the development-status checkboxes. Add `.github/ISSUE_TEMPLATE/bug_report.md` and `feature_request.md`. Keep to the project rule against naming or comparing to the reference implementation in committed docs.

**Blocked by:** 11, 09 (enough of the build to describe accurately); run last

**Status:** ready-for-agent

- [ ] `DEPENDENCIES.md` reflects the chosen frontend stack and drops the stale "future" entries; one line on why not shadcn
- [ ] Root `CLAUDE.md` frontend section matches the built stack
- [ ] `README.md` stack table and status list updated
- [ ] `.github/ISSUE_TEMPLATE/bug_report.md` and `feature_request.md` exist and match the Help-section links
- [ ] No committed doc names or compares against the reference implementation
- [ ] `CONTEXT.md` glossary re-read for accuracy against the final code; fix any drift
