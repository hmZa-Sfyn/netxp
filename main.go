package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	ModulesDir string            `json:"modules_dir"`
	Dirs       map[string]string `json:"dirs"`
	LastDir    string            `json:"last_dir"`
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
	cfg := &Config{Dirs: map[string]string{}}
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

func templateFor(lang string, name string) (string, error) {
	switch strings.ToLower(lang) {
	case "bash", "sh":
		return fmt.Sprintf("#!/usr/bin/env bash\n# %s - netxp module (bash)\necho \"Hello from %s (bash)\"\n", name, name), nil
	case "python", "py", "python3":
		return fmt.Sprintf("#!/usr/bin/env python3\n# %s - netxp module (python)\nprint(\"Hello from %s (python)\")\n", name, name), nil
	case "ruby", "rb":
		return fmt.Sprintf("#!/usr/bin/env ruby\n# %s - netxp module (ruby)\nputs 'Hello from %s (ruby)'\n", name, name), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
}

func listModules(cfg *Config) ([]string, error) {
	files, err := ioutil.ReadDir(cfg.ModulesDir)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		out = append(out, f.Name())
	}
	return out, nil
}

func runModule(cfg *Config, name string) error {
	// try to find file by prefix name
	files, _ := ioutil.ReadDir(cfg.ModulesDir)
	var target string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), name) || f.Name() == name {
			target = filepath.Join(cfg.ModulesDir, f.Name())
			break
		}
	}
	if target == "" {
		return fmt.Errorf("module not found: %s", name)
	}
	// make executable
	_ = os.Chmod(target, 0755)
	cmd := exec.Command(target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func execWithInterpreter(path, lang string) error {
	var cmd *exec.Cmd
	switch strings.ToLower(lang) {
	case "bash", "sh":
		cmd = exec.Command("bash", path)
	case "python", "py", "python3":
		python := "python3"
		if runtime.GOOS == "windows" {
			python = "python"
		}
		cmd = exec.Command(python, path)
	case "ruby", "rb":
		cmd = exec.Command("ruby", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("failed loading config:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("netxp - lightweight modular shell (type 'help')")
	for {
		fmt.Print("netxp> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]
		switch cmd {
		case "help":
			fmt.Println("commands: new <name> <lang>, run <name>, delete <name>, list, refresh, cd <path>, setdir <alias> <path>, gotodir <alias>, exit")
		case "new":
			if len(args) < 2 {
				fmt.Println("usage: new <name> <lang>")
				continue
			}
			name := args[0]
			lang := args[1]
			tmpl, e := templateFor(lang, name)
			if e != nil {
				fmt.Println(e)
				continue
			}
			ext := map[string]string{"bash": "sh", "sh": "sh", "python": "py", "python3": "py", "py": "py", "ruby": "rb", "rb": "rb"}[strings.ToLower(lang)]
			if ext == "" {
				ext = strings.ToLower(lang)
			}
			fname := filepath.Join(cfg.ModulesDir, fmt.Sprintf("%s.%s", name, ext))
			if err := ioutil.WriteFile(fname, []byte(tmpl), 0755); err != nil {
				fmt.Println("failed to create module:", err)
			} else {
				fmt.Println("created:", fname)
			}
		case "run":
			if len(args) < 1 {
				fmt.Println("usage: run <name>")
				continue
			}
			if err := runModule(cfg, args[0]); err != nil {
				fmt.Println("run error:", err)
			}
		case "delete":
			if len(args) < 1 {
				fmt.Println("usage: delete <name>")
				continue
			}
			// allow prefix match
			files, _ := ioutil.ReadDir(cfg.ModulesDir)
			found := false
			for _, f := range files {
				if f.IsDir() {
					continue
				}
				if strings.HasPrefix(f.Name(), args[0]) || f.Name() == args[0] {
					p := filepath.Join(cfg.ModulesDir, f.Name())
					_ = os.Remove(p)
					fmt.Println("deleted:", p)
					found = true
				}
			}
			if !found {
				fmt.Println("not found")
			}
		case "list":
			mods, _ := listModules(cfg)
			for _, m := range mods {
				fmt.Println(m)
			}
		case "refresh":
			fmt.Println("modules refreshed")
		case "cd":
			if len(args) < 1 {
				fmt.Println("usage: cd <path>")
				continue
			}
			p := args[0]
			if err := os.Chdir(p); err != nil {
				fmt.Println("cd failed:", err)
				continue
			}
			cfg.LastDir, _ = os.Getwd()
			_ = saveConfig(cfg)
			fmt.Println("cwd:", cfg.LastDir)
		case "setdir":
			if len(args) < 2 {
				fmt.Println("usage: setdir <alias> <path>")
				continue
			}
			alias := args[0]
			path := args[1]
			cfg.Dirs[alias] = path
			_ = saveConfig(cfg)
			fmt.Println("set dir", alias)
		case "gotodir":
			if len(args) < 1 {
				fmt.Println("usage: gotodir <alias>")
				continue
			}
			alias := args[0]
			p, ok := cfg.Dirs[alias]
			if !ok {
				fmt.Println("alias not found")
				continue
			}
			if err := os.Chdir(p); err != nil {
				fmt.Println("cd failed:", err)
				continue
			}
			cfg.LastDir, _ = os.Getwd()
			_ = saveConfig(cfg)
			fmt.Println("cwd:", cfg.LastDir)
		case "exit", "quit":
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command. type 'help'")
		}
	}
}
