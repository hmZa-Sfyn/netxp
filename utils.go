package main

import (
	"os/exec"
	"runtime"
)

func execModuleByExt(path, ext string) ([]byte, error) {
	// Choose interpreter by extension
	if ext == "sh" || ext == "bash" {
		return exec.Command("bash", path).CombinedOutput()
	}
	if ext == "py" || ext == "python3" {
		py := "python3"
		if runtime.GOOS == "windows" {
			py = "python"
		}
		return exec.Command(py, path).CombinedOutput()
	}
	if ext == "rb" || ext == "ruby" {
		return exec.Command("ruby", path).CombinedOutput()
	}
	return exec.Command(path).CombinedOutput()
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Colorize(s, code string) string {
	return code + s + CReset
}
