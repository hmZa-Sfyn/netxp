package moduling

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"netxp/config"
)

// Run executes a module and returns JSON output
func Run(cfg *config.Config, name string, args []string, input []byte) ([]byte, error) {
	files, err := ioutil.ReadDir(cfg.ModulesDir)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("module not found: %s", name)
	}

	_ = os.Chmod(target, 0755)
	cmd := exec.Command(target)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return []byte{}, nil
}

// Create creates a new module from template
func Create(cfg *config.Config, name, lang string) error {
	tmpl, err := template(lang, name)
	if err != nil {
		return err
	}
	ext := extensionFor(lang)
	fname := filepath.Join(cfg.ModulesDir, fmt.Sprintf("%s.%s", name, ext))
	return ioutil.WriteFile(fname, []byte(tmpl), 0755)
}

// Delete removes a module
func Delete(cfg *config.Config, name string) error {
	files, err := ioutil.ReadDir(cfg.ModulesDir)
	if err != nil {
		return err
	}
	found := false
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), name) || f.Name() == name {
			p := filepath.Join(cfg.ModulesDir, f.Name())
			if err := os.Remove(p); err != nil {
				return err
			}
			found = true
		}
	}
	if !found {
		return fmt.Errorf("module not found: %s", name)
	}
	return nil
}

// List returns all available modules
func List(cfg *config.Config) ([]string, error) {
	files, err := ioutil.ReadDir(cfg.ModulesDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, f.Name())
		}
	}
	return names, nil
}

func template(lang, name string) (string, error) {
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

func extensionFor(lang string) string {
	switch strings.ToLower(lang) {
	case "bash", "sh":
		return "sh"
	case "python", "py", "python3":
		return "py"
	case "ruby", "rb":
		return "rb"
	default:
		return strings.ToLower(lang)
	}
}
