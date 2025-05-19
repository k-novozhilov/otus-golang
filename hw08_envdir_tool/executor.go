package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
//
//nolint:gosec // запуск команд с пользовательскими аргументами разрешён заданием
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = buildEnv(env)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err == nil {
		if command.ProcessState != nil {
			return command.ProcessState.ExitCode()
		}
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}

func buildEnv(env Environment) []string {
	newEnv := os.Environ()
	for k, v := range env {
		if v.NeedRemove {
			for i := range newEnv {
				if stringIndex(newEnv[i], '=') == len(k) && newEnv[i][:len(k)] == k {
					newEnv[i] = ""
				}
			}
			continue
		}
		found := false
		for i := range newEnv {
			if stringIndex(newEnv[i], '=') == len(k) && newEnv[i][:len(k)] == k {
				newEnv[i] = k + "=" + v.Value
				found = true
				break
			}
		}
		if !found {
			newEnv = append(newEnv, k+"="+v.Value)
		}
	}
	filtered := make([]string, 0, len(newEnv))
	for _, e := range newEnv {
		if e != "" {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func stringIndex(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
