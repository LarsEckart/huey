package auth

import (
	"bufio"
	"cmp"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LarsEckart/huey/config"
	"github.com/LarsEckart/huey/hue"
)

// EnsureAuthenticated checks config and runs the auth flow if needed.
// Returns the loaded/updated config, or error if auth fails.
func EnsureAuthenticated() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Already configured?
	if cfg.IsConfigured() {
		return cfg, nil
	}

	// Need bridge IP?
	if cfg.BridgeIP == "" {
		ip, err := promptBridgeIP()
		if err != nil {
			return nil, err
		}
		cfg.BridgeIP = ip
	}

	// Need username?
	if cfg.Username == "" {
		username, err := registerWithBridge(cfg.BridgeIP)
		if err != nil {
			return nil, err
		}
		cfg.Username = username
	}

	// Save the config
	if err := cfg.Save(); err != nil {
		return nil, fmt.Errorf("save config: %w", err)
	}

	fmt.Println("✓ Configuration saved")
	return cfg, nil
}

// promptBridgeIP discovers the bridge IP address or asks the user for it.
func promptBridgeIP() (string, error) {
	fmt.Println("No Hue bridge configured.")
	fmt.Println("Discovering Hue bridges via https://discovery.meethue.com/ ...")

	bridges, err := discoverBridges()
	if err != nil {
		fmt.Printf("Automatic discovery failed: %v\n", err)
		return promptManualBridgeIP()
	}

	switch len(bridges) {
	case 0:
		fmt.Println("No Hue bridges found automatically.")
		return promptManualBridgeIP()
	case 1:
		fmt.Printf("✓ Found Hue bridge at %s\n", bridges[0].InternalIPAddress)
		return bridges[0].InternalIPAddress, nil
	default:
		return promptBridgeSelection(bridges)
	}
}

func discoverBridges() ([]hue.Bridge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	bridges, err := hue.DiscoverBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover bridges: %w", err)
	}
	return bridges, nil
}

func promptBridgeSelection(bridges []hue.Bridge) (string, error) {
	fmt.Println("Multiple Hue bridges found:")
	for i, bridge := range bridges {
		description := bridge.InternalIPAddress
		if bridge.ID != "" {
			description = fmt.Sprintf("%s (%s)", bridge.InternalIPAddress, bridge.ID)
		}
		fmt.Printf("  %d) %s\n", i+1, description)
	}

	input, err := promptString("\nEnter bridge number or IP address: ")
	if err != nil {
		return "", err
	}

	if selection, err := strconv.Atoi(input); err == nil {
		if selection < 1 || selection > len(bridges) {
			return "", fmt.Errorf("bridge selection must be between 1 and %d", len(bridges))
		}
		return bridges[selection-1].InternalIPAddress, nil
	}

	return input, nil
}

func promptManualBridgeIP() (string, error) {
	fmt.Println("Find your bridge IP in the Hue app: Settings → Hue Bridges → ⓘ")
	return promptString("\nEnter bridge IP address: ")
}

func promptString(prompt string) (string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	value := strings.TrimSpace(input)
	if value == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	return value, nil
}

// registerWithBridge prompts user to press link button, then registers.
func registerWithBridge(bridgeIP string) (string, error) {
	fmt.Println("\nTo authorize huey, press the link button on your Hue bridge.")
	fmt.Print("Press Enter when ready...")

	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	fmt.Println("Registering with bridge...")

	client := hue.NewClient(bridgeIP, "")
	hostname, _ := os.Hostname()
	deviceType := fmt.Sprintf("huey#%s", cmp.Or(hostname, "cli"))

	username, err := client.Register(deviceType)
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	fmt.Println("✓ Registered successfully")
	return username, nil
}
