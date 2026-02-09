package utils

import (
	"encoding/json"
	"strings"
)

// ParseCmd splits a command string into command name and arguments
func ParseCmd(line string) (string, []string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// SplitPipeline splits a pipeline string by '|'
func SplitPipeline(line string) []string {
	return strings.Split(line, "|")
}

// ParseJSON tries to parse bytes as JSON
func ParseJSON(data []byte) (interface{}, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)
	return v, err
}

// ToJSON converts value to JSON bytes
func ToJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return append(b, '\n')
}

// PrettyJSON returns indented JSON
func PrettyJSON(v interface{}) []byte {
	b, _ := json.MarshalIndent(v, "", "  ")
	return append(b, '\n')
}
