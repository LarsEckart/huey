## huey

A Go CLI app to control Philips Hue lights. Dual-mode: flag-based for scripting/agents, interactive TUI for humans.

Small steps, frequent commits.

## Documentation

- Official Hue API: https://developers.meethue.com/develop/get-started-2/
- Local API reference: [doc.md](doc.md)

## Tech Stack

- **Cobra** — CLI flag parsing
- **Bubble Tea** — TUI framework
- **Lip Gloss** — TUI styling

## Config

Stored in `~/.config/huey/config.json`:
- `bridge_ip` — Hue bridge IP address
- `username` — API username (obtained via bridge link button)

## First Run Flow

1. No config? → Prompt for bridge IP
2. No username? → "Press bridge button, then Enter" → POST to `/api`
3. Store both in config

## Commands

### Flag Mode (for agents/scripts)
```
huey lights              # List all lights
huey light <id> --on     # Turn light on
huey light <id> --off    # Turn light off
huey light <id> --toggle # Toggle light
```

### Interactive Mode (for humans)
```
huey                     # Opens TUI with light list, Enter to toggle
```

## API Endpoints Used

- `POST /api` — Register new username (requires bridge button press)
- `GET /api/{username}/lights` — List all lights
- `GET /api/{username}/lights/{id}` — Get light info
- `PUT /api/{username}/lights/{id}/state` — Change light state (`{"on": true/false}`)
