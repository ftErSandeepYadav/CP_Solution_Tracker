package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Platforms  map[string]*PlatformConfig `json:"platforms"`
	Repository RepositoryConfig           `json:"repository"`
	Sync       SyncConfig                 `json:"sync"`
}

type PlatformConfig struct {
	Enabled bool                   `json:"enabled"`
	Handle  string                 `json:"handle"`
	Cookies string                 `json:"cookies"`
	APIKey  string                 `json:"api_key,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

type RepositoryConfig struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Branch    string `json:"branch"`
	Token     string `json:"token"`
	LocalPath string `json:"local_path,omitempty"`
	UseGitHub bool   `json:"use_github"`
}

type SyncConfig struct {
	AutoSync           bool   `json:"auto_sync"`
	ConflictStrategy   string `json:"conflict_strategy"`
	StateFile          string `json:"state_file"`
	MaxConcurrency     int    `json:"max_concurrency"`
	OrganizedByContest bool   `json:"organized_by_contest"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func GetDefault() *Config {
	return &Config{
		Platforms: map[string]*PlatformConfig{
			"codeforces": &PlatformConfig{
				Enabled: true,
				Handle:  "Sandeep_Yadav_",
				Extra:   make(map[string]interface{}),
			},
		},
		Repository: RepositoryConfig{
			Owner:     "ftErSandeepYadav",
			Name:      "CP_Solutions",
			Branch:    "main",
			Token:     "",
			LocalPath: "",
			UseGitHub: true,
		},
		Sync: SyncConfig{
			AutoSync:           false,
			ConflictStrategy:   "keep_all",
			StateFile:          "sync_state.json",
			MaxConcurrency:     5,
			OrganizedByContest: true,
		},
	}
}

func (c *Config) GetPlatformConfig(name string) (*PlatformConfig, error) {
	cfg, exists := c.Platforms[name]
	if !exists {
		return nil, fmt.Errorf("platform '%s' not configured", name)
	}

	return cfg, nil
}
