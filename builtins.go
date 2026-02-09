package netxp
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"
)

func isBuiltin(name string) bool {
    switch name {
    case "pwd", "ls", "echo", "tab", "select", "new", "run", "list", "delete", "setdir", "gotodir", "workspaces":
        return true
    }
    return false
}

// runBuiltin executes internal commands and returns JSON bytes to be piped.
func runBuiltin(name string, args []string, input []byte, s *Shell) ([]byte, error) {
    switch name {
    case "pwd":
        cwd, _ := os.Getwd()
        j, _ := json.Marshal(map[string]string{"pwd": cwd})
        return append(j, '\n'), nil
    case "ls":
        path := "."
        if len(args) > 0 { path = args[0] }
        files, err := ioutil.ReadDir(path)
        if err != nil { return nil, err }
        out := []map[string]interface{}{}
        for _, f := range files {
            out = append(out, map[string]interface{}{
                "name": f.Name(),
                "size": f.Size(),
                "isdir": f.IsDir(),
                "mode": f.Mode().String(),
                "modtime": f.ModTime().Format(time.RFC3339),
            })
        }
        j, _ := json.Marshal(out)
        return append(j, '\n'), nil
    case "echo":
        // echoes parsed input JSON in a nicer format
        if len(input) == 0 && len(args) > 0 {
            return []byte(fmt.Sprintln(args...)), nil
        }
        if len(input) > 0 {
            // try to pretty print JSON
            var v interface{}
            if err := json.Unmarshal(input, &v); err == nil {
                pretty, _ := json.MarshalIndent(v, "", "  ")
                return append(pretty, '\n'), nil
            }
            return input, nil
        }
        return []byte("\n"), nil
    case "tab":
        // reads JSON array from input and returns pretty table as bytes
        if len(input) == 0 { return []byte("no input\n"), nil }
        // for simplicity just pretty-print JSON here
        var v interface{}
        if err := json.Unmarshal(input, &v); err != nil { return nil, err }
        pretty, _ := json.MarshalIndent(v, "", "  ")
        return append(pretty, '\n'), nil
    case "select":
        // args expected: fields comma separated
        if len(args) < 1 { return nil, fmt.Errorf("usage: select field1,field2") }
        if len(input) == 0 { return nil, fmt.Errorf("no input") }
        fields := args[0]
        fs := map[string]bool{}
        for _, f := range filepath.SplitList(fields) {
            fs[f] = true
        }
        var arr []map[string]interface{}
        if err := json.Unmarshal(input, &arr); err != nil { return nil, err }
        out := []map[string]interface{}{}
        for _, item := range arr {
            m := map[string]interface{}{}
            for k, v := range item {
                if fs[k] { m[k] = v }
            }
            out = append(out, m)
        }
        j, _ := json.Marshal(out)
        return append(j, '\n'), nil
    case "new", "run", "list", "delete", "setdir", "gotodir":
        // delegate to modules manager
        return runModuleBuiltin(name, args, s)
    case "workspaces":
        return []byte("[]\n"), nil
    }
    return nil, fmt.Errorf("unknown builtin: %s", name)
}
