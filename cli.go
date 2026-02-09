package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/peterh/liner"
)

type Shell struct {
	cfg   *Config
	repl  *liner.State
	histf string
}

func NewShell() *Shell {
	cfg, _ := loadConfig()
	repl := liner.NewLiner()
	repl.SetCtrlCAborts(true)
	histf := filepath.Join(configPath(), "history")
	if f, err := os.Open(histf); err == nil {
		repl.ReadHistory(f)
		f.Close()
	}
	return &Shell{cfg: cfg, repl: repl, histf: histf}
}

func (s *Shell) Close() {
	if s.repl != nil {
		if f, err := os.Create(s.histf); err == nil {
			s.repl.WriteHistory(f)
			f.Close()
		}
		s.repl.Close()
	}
}

func (s *Shell) Run() error {
	s.repl.SetMultiLineMode(false)
	fmt.Println(Colorize("netxp - modular shell", CInfo))
	for {
		line, err := s.repl.Prompt(fmt.Sprintf("%s> ", filepath.Base(s.cfg.ModulesDir)))
		if err == liner.ErrPromptAborted || err == io.EOF {
			fmt.Println()
			return nil
		}
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s.repl.AppendHistory(line)
		if line == "exit" || line == "quit" {
			return nil
		}
		if err := s.executePipeline(line); err != nil {
			s.PrettyError("error", err, "")
		}
	}
}

// executePipeline handles simple '|' pipelines between builtins and external commands.
func (s *Shell) executePipeline(line string) error {
	stages := splitPipeline(line)
	var input []byte
	var err error
	for i, stage := range stages {
		stage = strings.TrimSpace(stage)
		// parse command and args
		cmdName, args := parseCmd(stage)
		// builtin?
		if isBuiltin(cmdName) {
			input, err = runBuiltin(cmdName, args, input, s)
			if err != nil {
				return err
			}
			continue
		}
		// external command
		input, err = runExternal(stage, input)
		if err != nil {
			return fmt.Errorf("external '%s' failed: %w", stage, err)
		}
		// last stage: print if no next consumer
		if i == len(stages)-1 && len(input) > 0 {
			fmt.Print(string(input))
		}
	}
	return nil
}

func (s *Shell) PrettyError(kind string, err error, hint string) {
	fmt.Println(Colorize("error:", CErr), Colorize(err.Error(), CErr))
	if hint != "" {
		fmt.Println(Colorize("hint:", CWarn), hint)
	}
}

// parseCmd: very small tokenizer
func parseCmd(line string) (string, []string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func splitPipeline(line string) []string {
	// naive split by '|', does not support quoted pipes yet
	parts := strings.Split(line, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func runExternal(command string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("/bin/sh", "-c", command)
	}
	if len(input) > 0 {
		cmd.Stdin = bytes.NewReader(input)
	}
	out, err := cmd.CombinedOutput()
	return out, err
}
