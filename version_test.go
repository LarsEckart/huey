package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	cmd := newRootCmd()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(stdout.String(), "huey version ") {
		t.Fatalf("expected version output, got:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected empty stderr, got:\n%s", stderr.String())
	}
}

func TestAppVersionPrefersInjectedVersion(t *testing.T) {
	previous := version
	version = "v9.9.9"
	defer func() {
		version = previous
	}()

	if got := appVersion(); got != "v9.9.9" {
		t.Fatalf("expected injected version, got %q", got)
	}
}
