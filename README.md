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

- **↑/↓** or **j/k** — Navigate lights
- **Enter** or **Space** — Toggle selected light
- **q** — Quit

### Command Line

List all lights:
```bash
huey lights
```

Show light details:
```bash
huey light 1
```

Turn a light on/off:
```bash
huey light 1 --on
huey light 1 --off
```

Toggle a light:
```bash
huey light 1 --toggle
```

## Finding Your Bridge IP

- Check your router's connected devices
- Use the Hue app: Settings → Hue Bridges → (i) icon
- Visit https://discovery.meethue.com in your browser
