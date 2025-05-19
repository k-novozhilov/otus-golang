package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: go-envdir <envdir> <command> [args...]")
		os.Exit(1)
	}
	envdir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(envdir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := RunCmd(cmd, env)
	os.Exit(code)
}
