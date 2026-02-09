package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"netxp/builtins"
	"netxp/config"
	"netxp/moduling"
	"netxp/utils"

	"github.com/peterh/liner"
)

// Shell represents the interactive shell
type Shell struct {
	cfg   *config.Config
	repl  *liner.State
	histf string
}

// NewShell creates a new shell instance
func NewShell() *Shell {
	cfg, _ := config.Load()
	repl := liner.NewLiner()
	repl.SetCtrlCAborts(true)
	histf := config.HistoryFile()
	if f, err := os.Open(histf); err == nil {
		repl.ReadHistory(f)
		f.Close()
	}
	return &Shell{cfg: cfg, repl: repl, histf: histf}
}

// Close saves history and closes the shell
func (s *Shell) Close() {
	if s.repl != nil {
		if f, err := os.Create(s.histf); err == nil {
			s.repl.WriteHistory(f)
			f.Close()
		}
		s.repl.Close()
	}
}

// Run starts the interactive shell loop
func (s *Shell) Run() error {
	s.repl.SetMultiLineMode(false)
	fmt.Println(utils.Colorize("netxp - modular shell (type 'help')", utils.CInfo))
	for {
		prompt := fmt.Sprintf("%s> ", s.cfg.Workspace)
		if s.cfg.Workspace == "" {
			prompt = fmt.Sprintf("%s> ", filepath.Base(s.cfg.ModulesDir))
		}
		line, err := s.repl.Prompt(prompt)
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
			fmt.Println("bye")
			return nil
		}
		if line == "help" {
			s.printHelp()
			continue
		}
		if err := s.executePipeline(line); err != nil {
			s.PrettyError("error", err, "")
		}
	}
}

// executePipeline handles piped commands
func (s *Shell) executePipeline(line string) error {
	stages := utils.SplitPipeline(line)
	var input []byte
	var err error
	for i, stage := range stages {
		stage = strings.TrimSpace(stage)
		cmdName, args := utils.ParseCmd(stage)
		if cmdName == "" {
			continue
		}

		// Builtin command
		if builtins.IsBuiltin(cmdName) {
			input, err = builtins.Execute(cmdName, args, input)
			if err != nil {
				return err
			}
			if i == len(stages)-1 {
				fmt.Print(string(input))
			}
			continue
		}

		// Module command
		if strings.HasPrefix(cmdName, "run:") {
			modName := strings.TrimPrefix(cmdName, "run:")
			input, err = moduling.Run(s.cfg, modName, args, input)
			if err != nil {
				return err
			}
			if i == len(stages)-1 {
				fmt.Print(string(input))
			}
			continue
		}

		// External command
		cmd := exec.Command(cmdName, args...)
		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = os.Stderr
		cmd.Stdin = bytes.NewReader(input)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed: %s", err)
		}
		input = outBuf.Bytes()
		if i == len(stages)-1 {
			fmt.Print(string(input))
		}
	}
	return nil
}

// PrettyError prints a colorized error message
func (s *Shell) PrettyError(errType string, err error, hints string) {
	fmt.Printf("%s: %s\n", utils.ColorizeError(errType), err.Error())
	if hints != "" {
		fmt.Printf("  %s\n", utils.ColorizeWarn("hint: "+hints))
	}
}

// printHelp displays available commands
func (s *Shell) printHelp() {
	fmt.Println(utils.Colorize("\n=== NetXP Commands ===", utils.CBlue))
	fmt.Println("Module Commands:")
	fmt.Println("  new <name> <lang>     - Create new module (bash, python, ruby)")
	fmt.Println("  run <name> [args]     - Run a module")
	fmt.Println("  list                  - List all modules")
	fmt.Println("  delete <name>         - Delete a module")
	fmt.Println("\nDirectory Commands:")
	fmt.Println("  cd <path>             - Change directory")
	fmt.Println("  setdir <alias> <path> - Store directory alias")
	fmt.Println("  gotodir <alias>       - Go to stored directory")
	fmt.Println("  pwd                   - Print working directory")
	fmt.Println("\nBuiltin Commands:")
	builtins.InitDefaultBuiltins()
	for _, b := range builtins.List() {
		fmt.Printf("  %s\n", b)
	}
	fmt.Println("\nPiping:")
	fmt.Println("  cmd1 | cmd2 | cmd3    - Pipe JSON output between commands")
	fmt.Println("  cmd | tab             - Tabulate JSON output")
	fmt.Println("\nOther:")
	fmt.Println("  help                  - Show this help")
	fmt.Println("  exit, quit            - Exit shell")
	fmt.Println()
}
