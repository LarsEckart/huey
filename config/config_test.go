package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	// Use temp dir for test
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	_ = originalHome // t.Setenv auto-restores

	// Load should return empty config when file doesn't exist
	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if config.BridgeIP != "" || config.Username != "" {
		t.Errorf("Expected empty config, got: %+v", config)
	}
	if config.IsConfigured() {
		t.Error("Empty config should not be configured")
	}

	// Save config
	config.BridgeIP = "192.168.1.100"
	config.Username = "testuser123"
	if err := config.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists with correct permissions
	path := filepath.Join(tmpDir, ".config", "huey", "config.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Config file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected 0600 permissions, got %o", info.Mode().Perm())
	}

	// Load again and verify
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Load after save failed: %v", err)
	}
	if loadedConfig.BridgeIP != "192.168.1.100" {
		t.Errorf("BridgeIP mismatch: %s", loadedConfig.BridgeIP)
	}
	if loadedConfig.Username != "testuser123" {
		t.Errorf("Username mismatch: %s", loadedConfig.Username)
	}
	if !loadedConfig.IsConfigured() {
		t.Error("Filled config should be configured")
	}
}
