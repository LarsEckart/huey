package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

// promptBridgeIP asks the user for the bridge IP address.
func promptBridgeIP() (string, error) {
	fmt.Println("No Hue bridge configured.")
	fmt.Println("Find your bridge IP at: https://discovery.meethue.com/")
	fmt.Print("\nEnter bridge IP address: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	ip := strings.TrimSpace(input)
	if ip == "" {
		return "", fmt.Errorf("bridge IP cannot be empty")
	}

	return ip, nil
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
	if hostname == "" {
		hostname = "cli"
	}
	deviceType := fmt.Sprintf("huey#%s", hostname)

	username, err := client.Register(deviceType)
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	fmt.Println("✓ Registered successfully")
	return username, nil
}
