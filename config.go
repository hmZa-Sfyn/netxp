package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

type Config struct {
	ModulesDir string            `json:"modules_dir"`
	Dirs       map[string]string `json:"dirs"`
	LastDir    string            `json:"last_dir"`
	Theme      string            `json:"theme"`
}

func configPath() string {
	if runtime.GOOS == "windows" {
		if v := os.Getenv("APPDATA"); v != "" {
			return filepath.Join(v, "netxp")
		}
	}
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".netxp")
}

func loadConfig() (*Config, error) {
	cfgDir := configPath()
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return nil, err
	}
	cfgFile := filepath.Join(cfgDir, "config.json")
	cfg := &Config{Dirs: map[string]string{}, Theme: "default"}
	if b, err := ioutil.ReadFile(cfgFile); err == nil {
		_ = json.Unmarshal(b, cfg)
	}
	if cfg.ModulesDir == "" {
		cfg.ModulesDir = filepath.Join(cfgDir, "modules")
		_ = os.MkdirAll(cfg.ModulesDir, 0755)
	}
	return cfg, nil
}

func saveConfig(cfg *Config) error {
	cfgDir := configPath()
	cfgFile := filepath.Join(cfgDir, "config.json")
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return ioutil.WriteFile(cfgFile, b, 0644)
}
