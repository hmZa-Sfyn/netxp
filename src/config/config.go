package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Config holds netxp configuration
type Config struct {
	ModulesDir string            `json:"modules_dir"`
	Dirs       map[string]string `json:"dirs"`
	LastDir    string            `json:"last_dir"`
	Theme      string            `json:"theme"`
	Workspace  string            `json:"workspace"`
}

// ConfigPath returns the platform-specific config directory
func ConfigPath() string {
	if runtime.GOOS == "windows" {
		if v := os.Getenv("APPDATA"); v != "" {
			return filepath.Join(v, "netxp")
		}
	}
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".netxp")
}

// Load reads config from disk
func Load() (*Config, error) {
	cfgDir := ConfigPath()
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return nil, err
	}
	cfgFile := filepath.Join(cfgDir, "config.json")
	cfg := &Config{Dirs: make(map[string]string)}
	if b, err := ioutil.ReadFile(cfgFile); err == nil {
		_ = json.Unmarshal(b, cfg)
	}
	if cfg.ModulesDir == "" {
		cfg.ModulesDir = filepath.Join(cfgDir, "modules")
		_ = os.MkdirAll(cfg.ModulesDir, 0755)
	}
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}
	return cfg, nil
}

// Save writes config to disk
func (c *Config) Save() error {
	cfgDir := ConfigPath()
	cfgFile := filepath.Join(cfgDir, "config.json")
	b, _ := json.MarshalIndent(c, "", "  ")
	return ioutil.WriteFile(cfgFile, b, 0644)
}

// HistoryFile returns path to history file
func HistoryFile() string {
	return filepath.Join(ConfigPath(), "history")
}

// WorkspacesDir returns path to workspaces directory
func WorkspacesDir() string {
	return filepath.Join(ConfigPath(), "workspaces")
}
