## v1 Scope

- [x] Project setup (Go module, dependencies)
- [x] Config load/save
- [x] HTTP client for bridge API
- [x] Auth flow (bridge button registration)
- [x] Flag mode commands
- [x] Interactive TUI (light list, toggle)

## v2 Scope

- [x] Groups (list, toggle)
- [x] Light/Group renaming (CLI --name flag, TUI 'r' key)
- [ ] Scenes (list, activate)

## Backlog

- `--version` flag using `runtime/debug.BuildInfo` (auto-populated by `go install @tag`)
- `--json` output for scripting/agent use
- Brightness/color control
- Bridge discovery (mDNS)
- Room-aware views
- Favorites/presets (user-defined states)
- Scheduling
- Multi-bridge support
