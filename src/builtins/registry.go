package builtins

import (
	"encoding/json"
	"fmt"
)

// ExecutionError represents a structured error from a builtin command
type ExecutionError struct {
	Command string      `json:"command"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Hints   []string    `json:"hints,omitempty"`
	Context interface{} `json:"context,omitempty"`
}

// BuiltinFunc is the signature for a builtin command handler
// It receives: command name, args, input (from pipe), and returns output bytes and error
type BuiltinFunc func(name string, args []string, input []byte) ([]byte, error)

// Registry holds all registered builtin commands
var Registry = make(map[string]BuiltinFunc)

// Register adds a builtin command to the registry
func Register(name string, fn BuiltinFunc) {
	Registry[name] = fn
}

// Execute runs a builtin command by name
func Execute(name string, args []string, input []byte) ([]byte, error) {
	fn, exists := Registry[name]
	if !exists {
		return nil, fmt.Errorf("command not found: %s", name)
	}
	return fn(name, args, input)
}

// IsBuiltin checks if a command is registered
func IsBuiltin(name string) bool {
	_, exists := Registry[name]
	return exists
}

// List returns all registered builtins
func List() []string {
	names := []string{}
	for name := range Registry {
		names = append(names, name)
	}
	return names
}

// StructuredError creates a JSON error response
func StructuredError(cmd string, code int, msg string, hints []string) []byte {
	err := ExecutionError{
		Command: cmd,
		Code:    code,
		Message: msg,
		Hints:   hints,
	}
	b, _ := json.Marshal(err)
	return append(b, '\n')
}

// StructuredOutput wraps data as JSON with success flag
func StructuredOutput(data interface{}) []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"success": true,
		"data":    data,
	})
	return append(b, '\n')
}
