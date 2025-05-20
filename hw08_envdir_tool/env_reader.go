package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.Contains(name, "=") {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			env[name] = EnvValue{"", true}
			continue
		}
		lineEnd := bytes.IndexByte(data, '\n')
		var line []byte
		if lineEnd == -1 {
			line = data
		} else {
			line = data[:lineEnd]
		}
		line = bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})
		str := strings.TrimRight(string(line), " \t")
		env[name] = EnvValue{str, false}
	}
	return env, nil
}
