## huey

A Go CLI app to control Philips Hue lights. Dual-mode: flag-based for scripting/agents, interactive TUI for humans.

Small steps, frequent commits.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│   CLI (cmd/*.go)              TUI (tui/tui.go)              │
│   - fmt.Printf output         - Bubble Tea UI               │
│   - Single operation          - Interactive session         │
│        │                            │                       │
│        └──────────┬─────────────────┘                       │
│                   ▼                                         │
│            hue/client.go                                    │
│            (HTTP calls to bridge)                           │
└─────────────────────────────────────────────────────────────┘
```

**All features must be implemented for both CLI and TUI.**

## Tech Stack

- **Cobra** — CLI flag parsing
- **Bubble Tea** — TUI framework
- **Lip Gloss** — TUI styling

## Hue Bridge API

Base URL: `http://<bridge_ip>/api`

### Authentication
- `POST /api` — Register username (requires bridge button press)
  - Body: `{"devicetype": "app#device"}`
  - Returns: `{"success": {"username": "..."}}`

### Lights
- `GET /api/{username}/lights` — List all lights (returns map of ID → light)
- `GET /api/{username}/lights/{id}` — Get single light
- `PUT /api/{username}/lights/{id}/state` — Set light state
  - Body: `{"on": bool, "bri": 0-254, "hue": 0-65535, "sat": 0-254}`
- `PUT /api/{username}/lights/{id}` — Update light attributes
  - Body: `{"name": "..."}`

### Groups
Groups have types: `Room` (light can be in one), `Zone` (light can be in many), `Entertainment`, `LightGroup`

- `GET /api/{username}/groups` — List all groups
- `POST /api/{username}/groups` — Create group
  - Body: `{"name": "...", "type": "Room"|"Zone", "lights": ["1","2"]}`
- `PUT /api/{username}/groups/{id}` — Update group attributes
  - Body: `{"name": "..."}`
- `PUT /api/{username}/groups/{id}/action` — Set state for all lights in group
  - Body: `{"on": bool}` or `{"scene": "sceneId"}`
- `DELETE /api/{username}/groups/{id}` — Delete group

### Scenes
Scenes are saved light configurations for a group.

- `GET /api/{username}/scenes` — List all scenes
- `GET /api/{username}/scenes/{id}` — Get scene details
- `POST /api/{username}/scenes` — Create scene (captures current light states)
  - Body: `{"name": "...", "type": "GroupScene", "group": "groupId", "recycle": false}`
  - Returns: `{"success": {"id": "..."}}`
- `PUT /api/{username}/groups/{groupId}/action` — Activate scene
  - Body: `{"scene": "sceneId"}`
- `DELETE /api/{username}/scenes/{id}` — Delete scene
