package main

import (
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := "testdata/env"
	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	expect := map[string]EnvValue{
		"HELLO": {"\"hello\"", false},
		"BAR":   {"bar", false},
		"FOO":   {"   foo\nwith new line", false},
		"UNSET": {"", true},
		"EMPTY": {"", false},
	}
	for k, v := range expect {
		got, ok := env[k]
		if !ok {
			t.Errorf("missing key %s", k)
			continue
		}
		if got != v {
			t.Errorf("%s: got %+v, want %+v", k, got, v)
		}
	}
	for k := range env {
		if _, ok := expect[k]; !ok {
			t.Errorf("unexpected key %s", k)
		}
	}
}
