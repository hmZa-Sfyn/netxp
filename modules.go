package netxp
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

func runModuleBuiltin(name string, args []string, s *Shell) ([]byte, error) {
    switch name {
    case "new":
        if len(args) < 2 { return nil, fmt.Errorf("usage: new <name> <lang>") }
        nm := args[0]
        lang := args[1]
        tmpl, err := templateFor(lang, nm)
        if err != nil { return nil, err }
        ext := map[string]string{"bash":"sh","sh":"sh","python":"py","python3":"py","py":"py","ruby":"rb","rb":"rb"}[strings.ToLower(lang)]
        if ext == "" { ext = strings.ToLower(lang) }
        fname := filepath.Join(s.cfg.ModulesDir, fmt.Sprintf("%s.%s", nm, ext))
        if err := ioutil.WriteFile(fname, []byte(tmpl), 0755); err != nil { return nil, err }
        out, _ := json.Marshal(map[string]string{"created": fname})
        return append(out, '\n'), nil
    case "list":
        files, _ := ioutil.ReadDir(s.cfg.ModulesDir)
        arr := []map[string]interface{}{}
        for _, f := range files {
            if f.IsDir() { continue }
            arr = append(arr, map[string]interface{}{"name": f.Name(), "size": f.Size()})
        }
        j, _ := json.Marshal(arr)
        return append(j, '\n'), nil
    case "delete":
        if len(args) < 1 { return nil, fmt.Errorf("usage: delete <name>") }
        files, _ := ioutil.ReadDir(s.cfg.ModulesDir)
        found := false
        for _, f := range files {
            if strings.HasPrefix(f.Name(), args[0]) {
                _ = os.Remove(filepath.Join(s.cfg.ModulesDir, f.Name()))
                found = true
            }
        }
        if !found { return nil, fmt.Errorf("not found") }
        return []byte("{}\n"), nil
    case "run":
        if len(args) < 1 { return nil, fmt.Errorf("usage: run <name>") }
        // find file by prefix
        files, _ := ioutil.ReadDir(s.cfg.ModulesDir)
        var target string
        for _, f := range files {
            if strings.HasPrefix(f.Name(), args[0]) { target = filepath.Join(s.cfg.ModulesDir, f.Name()); break }
        }
        if target == "" { return nil, fmt.Errorf("module not found") }
        data, err := ioutil.ReadFile(target)
        if err != nil { return nil, err }
        // attempt to execute with interpreter based on extension
        ext := filepath.Ext(target)
        if len(ext) > 0 { ext = ext[1:] }
        // write to temp file and execute
        out, err := execModuleByExt(target, ext)
        return out, err
    case "setdir":
        if len(args) < 2 { return nil, fmt.Errorf("usage: setdir <alias> <path>") }
        s.cfg.Dirs[args[0]] = args[1]
        _ = saveConfig(s.cfg)
        out, _ := json.Marshal(map[string]string{"set": args[0]})
        return append(out, '\n'), nil
    case "gotodir":
        if len(args) < 1 { return nil, fmt.Errorf("usage: gotodir <alias>") }
        p, ok := s.cfg.Dirs[args[0]]
        if !ok { return nil, fmt.Errorf("alias not found") }
        if err := os.Chdir(p); err != nil { return nil, err }
        s.cfg.LastDir, _ = os.Getwd()
        _ = saveConfig(s.cfg)
        out, _ := json.Marshal(map[string]string{"cwd": s.cfg.LastDir})
        return append(out, '\n'), nil
    }
    return nil, fmt.Errorf("unknown module command")
}
