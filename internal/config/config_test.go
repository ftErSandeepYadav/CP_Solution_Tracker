package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_LoadSave(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create config
	cfg := GetDefault()
	cfg.Platforms["codeforces"].Handle = "testuser"
	cfg.Repository.Owner = "testowner"
	cfg.Repository.Name = "test-repo"
	cfg.Repository.Token = "ghp_test123"

	// Save
	err := cfg.Save(configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify data integrity
	if loaded.Platforms["codeforces"].Handle != "testuser" {
		t.Errorf("Expected handle 'testuser', got '%s'",
			loaded.Platforms["codeforces"].Handle)
	}

	if loaded.Repository.Owner != "testowner" {
		t.Errorf("Expected owner 'testowner', got '%s'",
			loaded.Repository.Owner)
	}

	if loaded.Repository.Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got '%s'",
			loaded.Repository.Name)
	}
}

func TestConfig_LoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}
}

func TestGetDefault(t *testing.T) {
	cfg := GetDefault()

	// Check defaults
	if cfg.Repository.Branch != "main" {
		t.Errorf("Expected default branch 'main', got '%s'",
			cfg.Repository.Branch)
	}

	if !cfg.Platforms["codeforces"].Enabled {
		t.Error("Expected codeforces to be enabled by default")
	}

	if cfg.Sync.ConflictStrategy != "keep_all" {
		t.Errorf("Expected default conflict strategy 'keep_all', got '%s'",
			cfg.Sync.ConflictStrategy)
	}

	if cfg.Sync.MaxConcurrency != 5 {
		t.Errorf("Expected default max concurrency 5, got %d",
			cfg.Sync.MaxConcurrency)
	}
}

func TestConfig_GetPlatformConfig(t *testing.T) {
	cfg := GetDefault()

	// Get existing platform
	cfCfg, err := cfg.GetPlatformConfig("codeforces")
	if err != nil {
		t.Fatalf("Failed to get codeforces config: %v", err)
	}

	if !cfCfg.Enabled {
		t.Error("Expected codeforces to be enabled")
	}

	// Get non-existent platform
	_, err = cfg.GetPlatformConfig("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent platform")
	}
}

func TestValidator_ValidConfig(t *testing.T) {
	cfg := GetDefault()
	cfg.Platforms["codeforces"].Handle = "testuser"
	cfg.Repository.Owner = "testowner"
	cfg.Repository.Name = "test-repo"
	cfg.Repository.Token = "ghp_test123"

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestValidator_MissingHandle(t *testing.T) {
	cfg := GetDefault()
	// Handle is empty but platform is enabled
	cfg.Platforms["codeforces"].Handle = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing handle")
	}

	verr, ok := err.(ValidationErrors)
	if !ok {
		t.Fatal("Expected ValidationErrors type")
	}

	found := false
	for _, e := range verr {
		if e.Field == "platforms.codeforces.handle" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error for platforms.codeforces.handle")
	}
}

func TestValidator_MissingGitHubConfig(t *testing.T) {
	cfg := GetDefault()
	cfg.Platforms["codeforces"].Handle = "testuser"
	// Repository fields are empty but UseGitHub is true
	cfg.Repository.Owner = ""
	cfg.Repository.Name = ""
	cfg.Repository.Token = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing GitHub config")
	}

	verr, ok := err.(ValidationErrors)
	if !ok {
		t.Fatal("Expected ValidationErrors type")
	}

	// Should have errors for owner, name, and token
	if len(verr) < 3 {
		t.Errorf("Expected at least 3 validation errors, got %d", len(verr))
	}
}

func TestValidator_InvalidConflictStrategy(t *testing.T) {
	cfg := GetDefault()
	cfg.Platforms["codeforces"].Handle = "testuser"
	cfg.Repository.Owner = "testowner"
	cfg.Repository.Name = "test-repo"
	cfg.Repository.Token = "ghp_test123"
	cfg.Sync.ConflictStrategy = "invalid_strategy"

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid conflict strategy")
	}

	verr, ok := err.(ValidationErrors)
	if !ok {
		t.Fatal("Expected ValidationErrors type")
	}

	found := false
	for _, e := range verr {
		if e.Field == "sync.conflict_strategy" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error for sync.conflict_strategy")
	}
}

func TestValidator_MaxConcurrency(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
		shouldError bool
	}{
		{"zero concurrency", 0, true},
		{"valid concurrency", 5, false},
		{"max valid", 20, false},
		{"too high", 25, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := GetDefault()
			cfg.Platforms["codeforces"].Handle = "someHandle"
			cfg.Repository.Owner = "owner"
			cfg.Repository.Name = "repo"
			cfg.Repository.Token = "token"
			cfg.Sync.MaxConcurrency = tt.concurrency

			err := cfg.Validate()
			if tt.shouldError && err == nil {
				t.Error("Expected validation error")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
