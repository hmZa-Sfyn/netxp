package builtins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// InitDefaultBuiltins registers all default system builtins
func InitDefaultBuiltins() {
	Register("pwd", CmdPwd)
	Register("ls", CmdLs)
	Register("echo", CmdEcho)
	Register("tab", CmdTab)
	Register("select", CmdSelect)
	Register("cat", CmdCat)
	Register("cd", CmdCd)
	Register("env", CmdEnv)
	Register("whoami", CmdWhoami)
	Register("date", CmdDate)
	Register("mkdir", CmdMkdir)
	Register("rm", CmdRm)
	Register("cp", CmdCp)
	Register("mv", CmdMv)
	Register("find", CmdFind)
	Register("grep", CmdGrep)
	Register("wc", CmdWc)
	Register("head", CmdHead)
	Register("tail", CmdTail)
}

// CmdPwd returns current working directory
func CmdPwd(name string, args []string, input []byte) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return StructuredError(name, 1, err.Error(), []string{"ensure you have read permissions on current directory"}), nil
	}
	return StructuredOutput(map[string]string{"pwd": cwd}), nil
}

// CmdLs lists directory contents
func CmdLs(name string, args []string, input []byte) ([]byte, error) {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return StructuredError(name, 1, err.Error(), []string{"path not found or not accessible"}), nil
	}
	out := []map[string]interface{}{}
	for _, f := range files {
		out = append(out, map[string]interface{}{
			"name":    f.Name(),
			"size":    f.Size(),
			"isdir":   f.IsDir(),
			"mode":    f.Mode().String(),
			"modtime": f.ModTime().Format(time.RFC3339),
		})
	}
	return StructuredOutput(out), nil
}

// CmdEcho echoes or pretty-prints JSON input
func CmdEcho(name string, args []string, input []byte) ([]byte, error) {
	if len(input) == 0 && len(args) > 0 {
		return []byte(fmt.Sprintln(strings.Join(args, " "))), nil
	}
	if len(input) > 0 {
		var v interface{}
		if err := json.Unmarshal(input, &v); err == nil {
			pretty, _ := json.MarshalIndent(v, "", "  ")
			return append(pretty, '\n'), nil
		}
		return input, nil
	}
	return []byte("\n"), nil
}

// CmdTab formats JSON as table (simplified)
func CmdTab(name string, args []string, input []byte) ([]byte, error) {
	if len(input) == 0 {
		return StructuredError(name, 1, "no input", []string{"pipe data to tab command"}), nil
	}
	var v interface{}
	if err := json.Unmarshal(input, &v); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"input must be valid JSON"}), nil
	}
	pretty, _ := json.MarshalIndent(v, "", "  ")
	return append(pretty, '\n'), nil
}

// CmdSelect filters JSON fields
func CmdSelect(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing fields argument", []string{"usage: select field1,field2,field3"}), nil
	}
	if len(input) == 0 {
		return StructuredError(name, 1, "no input", []string{"pipe data to select"}), nil
	}
	fields := strings.Split(args[0], ",")
	fs := make(map[string]bool)
	for _, f := range fields {
		fs[strings.TrimSpace(f)] = true
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(input, &arr); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"input must be JSON array"}), nil
	}
	out := []map[string]interface{}{}
	for _, item := range arr {
		row := make(map[string]interface{})
		for k, v := range item {
			if fs[k] {
				row[k] = v
			}
		}
		if len(row) > 0 {
			out = append(out, row)
		}
	}
	return StructuredOutput(out), nil
}

// CmdCat reads file contents
func CmdCat(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing file argument", []string{"usage: cat <file>"}), nil
	}
	content, err := ioutil.ReadFile(args[0])
	if err != nil {
		return StructuredError(name, 1, err.Error(), []string{"file not found or not readable"}), nil
	}
	return content, nil
}

// CmdCd changes directory
func CmdCd(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing path argument", []string{"usage: cd <path>"}), nil
	}
	if err := os.Chdir(args[0]); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"directory not found or not accessible"}), nil
	}
	cwd, _ := os.Getwd()
	return StructuredOutput(map[string]string{"pwd": cwd}), nil
}

// CmdEnv lists environment variables
func CmdEnv(name string, args []string, input []byte) ([]byte, error) {
	env := os.Environ()
	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return StructuredOutput(envMap), nil
}

// CmdWhoami returns current user
func CmdWhoami(name string, args []string, input []byte) ([]byte, error) {
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}
	return StructuredOutput(map[string]string{"user": user}), nil
}

// CmdDate returns current time
func CmdDate(name string, args []string, input []byte) ([]byte, error) {
	return StructuredOutput(map[string]string{"date": time.Now().Format(time.RFC3339)}), nil
}

// CmdMkdir creates directory
func CmdMkdir(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing path", []string{"usage: mkdir <path>"}), nil
	}
	if err := os.MkdirAll(args[0], 0755); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"failed to create directory"}), nil
	}
	return StructuredOutput(map[string]string{"created": args[0]}), nil
}

// CmdRm removes file
func CmdRm(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing path", []string{"usage: rm <path>"}), nil
	}
	if err := os.RemoveAll(args[0]); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"failed to remove"}), nil
	}
	return StructuredOutput(map[string]string{"removed": args[0]}), nil
}

// CmdCp copies file
func CmdCp(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 2 {
		return StructuredError(name, 1, "missing src/dst", []string{"usage: cp <src> <dst>"}), nil
	}
	cmd := exec.Command("cp", "-r", args[0], args[1])
	if err := cmd.Run(); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"copy failed"}), nil
	}
	return StructuredOutput(map[string]string{"copied": args[0] + " -> " + args[1]}), nil
}

// CmdMv moves/renames file
func CmdMv(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 2 {
		return StructuredError(name, 1, "missing src/dst", []string{"usage: mv <src> <dst>"}), nil
	}
	if err := os.Rename(args[0], args[1]); err != nil {
		return StructuredError(name, 1, err.Error(), []string{"move failed"}), nil
	}
	return StructuredOutput(map[string]string{"moved": args[0] + " -> " + args[1]}), nil
}

// CmdFind searches for files
func CmdFind(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing pattern", []string{"usage: find <pattern>"}), nil
	}
	var matches []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.Contains(info.Name(), args[0]) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return StructuredError(name, 1, err.Error(), []string{"find failed"}), nil
	}
	return StructuredOutput(matches), nil
}

// CmdGrep searches in text
func CmdGrep(name string, args []string, input []byte) ([]byte, error) {
	if len(args) < 1 {
		return StructuredError(name, 1, "missing pattern", []string{"usage: grep <pattern>"}), nil
	}
	lines := strings.Split(string(input), "\n")
	var matches []string
	for _, line := range lines {
		if strings.Contains(line, args[0]) {
			matches = append(matches, line)
		}
	}
	return StructuredOutput(matches), nil
}

// CmdWc counts words/lines
func CmdWc(name string, args []string, input []byte) ([]byte, error) {
	text := string(input)
	lines := len(strings.Split(text, "\n"))
	words := len(strings.Fields(text))
	chars := len(text)
	return StructuredOutput(map[string]int{
		"lines": lines,
		"words": words,
		"chars": chars,
	}), nil
}

// CmdHead returns first lines
func CmdHead(name string, args []string, input []byte) ([]byte, error) {
	count := 10
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &count)
	}
	lines := strings.Split(string(input), "\n")
	if len(lines) > count {
		lines = lines[:count]
	}
	return StructuredOutput(lines), nil
}

// CmdTail returns last lines
func CmdTail(name string, args []string, input []byte) ([]byte, error) {
	count := 10
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &count)
	}
	lines := strings.Split(string(input), "\n")
	start := len(lines) - count
	if start < 0 {
		start = 0
	}
	return StructuredOutput(lines[start:]), nil
}
