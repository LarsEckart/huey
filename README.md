# huey

A simple CLI to control your Philips Hue lights.

## Installation

```bash
go install github.com/LarsEckart/huey@latest
```

Or download a binary from the [releases page](https://github.com/lars/huey/releases).

## First Run

1. Run `huey`
2. Enter your Hue bridge IP address when prompted
3. Press the link button on your Hue bridge
4. Press Enter

That's it! Your credentials are saved to `~/.config/huey/config.json`.

## Usage

### Interactive Mode

Run without arguments to open the interactive interface:

```bash
huey
```

**Navigation:**
- **Tab** or **l/h** — Switch between Lights, Groups, and Scenes tabs
- **↑/↓** or **j/k** — Navigate list
- **Space** — Toggle selected light/group, or activate scene
- **r** — Rename selected item (Lights/Groups only)
- **q** — Quit

**Groups tab only:**
- **a** — Add new group (room or zone)
- **i** — Show group info (lights in group)
- **d** — Delete group (with confirmation)

**Scenes tab only:**
- **a** — Add new scene (captures current light states)

### Command Line

#### Lights

List all lights:
```bash
huey lights
```

Show light details:
```bash
huey light 1
```

Control a light:
```bash
huey light 1 --on
huey light 1 --off
huey light 1 --toggle
```

Rename a light:
```bash
huey light 1 --name "Desk Lamp"
```

#### Groups

List all groups:
```bash
huey groups
```

Show group details:
```bash
huey group 1
```

Control a group:
```bash
huey group 1 --on
huey group 1 --off
huey group 1 --toggle
```

Rename a group:
```bash
huey group 1 --name "Living Room"
```

Delete a group:
```bash
huey group 1 --delete
```

Create a group:
```bash
huey group-create --name "My Zone" --type zone --lights 1,2,3
huey group-create --name "My Room" --type room
```

#### Scenes

List all scenes:
```bash
huey scenes
```

Activate a scene:
```bash
huey scene TqkSoVMtx4juUbU
```

Create a scene (captures current light states):
```bash
huey scene-create --name "Focus" --group 85
```

## Finding Your Bridge IP

- Check your router's connected devices
- Use the Hue app: Settings → Hue Bridges → (i) icon
- Visit https://discovery.meethue.com in your browser
