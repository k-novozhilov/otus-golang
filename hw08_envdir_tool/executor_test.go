package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO":   {"bar", false},
		"HELLO": {"world", false},
		"UNSET": {"", true},
	}
	cmd := []string{"bash", "-c", "echo $FOO $HELLO $UNSET"}
	var out bytes.Buffer
	command := exec.Command(cmd[0], cmd[1:]...)
	osenv := os.Environ()
	for k, v := range env {
		if v.NeedRemove {
			for i := range osenv {
				if strings.HasPrefix(osenv[i], k+"=") {
					osenv[i] = ""
				}
			}
			continue
		}
		found := false
		for i := range osenv {
			if strings.HasPrefix(osenv[i], k+"=") {
				osenv[i] = k + "=" + v.Value
				found = true
				break
			}
		}
		if !found {
			osenv = append(osenv, k+"="+v.Value)
		}
	}
	filtered := make([]string, 0, len(osenv))
	for _, e := range osenv {
		if e != "" {
			filtered = append(filtered, e)
		}
	}
	command.Env = filtered
	command.Stdout = &out
	command.Stderr = &out
	code := 0
	err := command.Run()
	var exitErr *exec.ExitError
	switch {
	case err != nil && errors.As(err, &exitErr):
		code = exitErr.ExitCode()
	case err != nil && command.ProcessState != nil:
		code = command.ProcessState.ExitCode()
	case err != nil:
		t.Fatalf("run error: %v", err)
	case command.ProcessState != nil:
		code = command.ProcessState.ExitCode()
	}
	if code != 0 {
		t.Errorf("exit code: got %d, want 0", code)
	}
	res := strings.TrimSpace(out.String())
	if res != "bar world" {
		t.Errorf("output: got %q, want %q", res, "bar world")
	}
}
