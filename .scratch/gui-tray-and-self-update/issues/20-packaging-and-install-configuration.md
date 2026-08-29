# 20: Packaging and install configuration

**What to build:** `task package` produced `scaffold-amd64-installer.exe`, and installing it created a program called "My Product" from "My Company". `build/config.yml` was filled in correctly but nothing ever consumed it: Wails compiles that file into `build/windows/info.json` and `build/windows/nsis/wails_tools.nsh`, and `wails3 update build-assets` had never been run. Configure the whole packaging path so the installer, the installed program, and the release artifacts all carry the real identity.

**Blocked by:** 09, 10

**Status:** done

- [x] `build/config.yml` carries the real identity: company "Haze", product "League RPC", identifier `com.its-haze.league-rpc` (matching `singleInstanceID`, which the installer's running-instance check names)
- [x] `task assets:update` regenerates the three Windows files into a scratch dir and copies only those back, so a Windows-only repo does not grow `darwin/`, `linux/`, `ios/` and an `nfpm.yaml`
- [x] Shipped binary renamed `league-rpc.exe`; the Go package stays at `cmd/league-rpc-gui` so the generated TS bindings keep their import path
- [x] `.syso` is written into the main package's directory. It was going to the repo root, where the linker never looks, so the built `.exe` had no icon and a blank Details tab
- [x] `info.json` is retargeted from language `0000` to `0409` after generation, and gains `product_version`. Under `0000` the Details tab reads back empty
- [x] A release tag stamps both the installer (`-DINFO_PRODUCTVERSION`) and the executable (a temporary `info.json` for the syso), trimmed from `v1.4.2-rc3` to the 3-part `1.4.2` NSIS wants
- [x] Per-user install scope by default: `$LOCALAPPDATA\Programs\League RPC`, no UAC. See ADR 0006
- [x] Uninstall key is `LeagueRPC` rather than the `HazeLeague RPC` concatenation
- [x] Components page with the desktop shortcut optional and unchecked; Start Menu shortcut always; finish-page "Start League RPC", defined only for user scope so an elevated installer never hands the app an admin token
- [x] Install and uninstall both abort with a Retry/Cancel prompt while the app is running, detected through the single-instance mutex
- [x] Uninstall deletes the `HKCU\...\Run` value `LeagueRPC`; settings and logs survive
- [x] Release publishes the installer alongside the binary, both covered by `SHA256SUMS`; `cmd/sign-release` takes a repeatable `-artifact`

**The updater bug this turned up:** `github.DefaultAssetMatcher` picks the first
asset whose lowercased name contains both the platform and the arch, defaulting
to the running build's. Releases published `league-rpc-gui.exe`, which contains
neither, so the matcher returned -1 and self-update would have found no asset on
every release. Ticket 10's dry run checked the signature sidecar by hand and
never exercised the download. The published name is now
`internal/updates.ReleaseAsset`, and `releaseasset_test.go` asserts the matcher
selects it, that it matches the running platform, and that `release.yml` still
names it. The matcher skips anything containing `-installer.`, so attaching the
installer to the same release is safe by construction.

**Verifying embedded metadata:** `[Diagnostics.FileVersionInfo]` and
`Get-Item .VersionInfo` report every string field blank for a version resource
that `VerQueryValueW` and Explorer's Details tab both read correctly. Check with
`Shell.Application`'s `GetDetailsOf` instead.

**Follow-ups:** `internal/updates.RepoSlug` still points at the test remote and
needs one line changed at repo migration. Code signing and a `LICENSE` file are
deliberately still deferred; the `signtool` hooks in `project.nsi` are left
commented for that work.
