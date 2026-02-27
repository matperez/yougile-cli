package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigPathCmd_PrintsResolvedPath(t *testing.T) {
	wantPath := filepath.Join(t.TempDir(), "config.yaml")
	resolvePath := func() (string, error) { return wantPath, nil }

	c := NewConfigPathCmd(resolvePath)
	c.SetOut(new(bytes.Buffer))
	c.SetErr(new(bytes.Buffer))

	err := c.Execute()
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := c.OutOrStdout().(*bytes.Buffer).String()
	got := strings.TrimSpace(out)
	if got != wantPath {
		t.Errorf("output = %q, want %q", got, wantPath)
	}
}

func TestConfigShowCmd_HumanOutput_MasksAPIKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(path, []byte("base_url: https://x.com\napi_key: secret123\n"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	resolvePath := func() (string, error) { return path, nil }
	outputJSON := func() bool { return false }

	c := NewConfigShowCmd(resolvePath, outputJSON)
	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(new(bytes.Buffer))

	err = c.Execute()
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "secret123") {
		t.Error("human output must not contain raw api_key")
	}
	if !strings.Contains(out, "***") {
		t.Error("human output must mask api_key with ***")
	}
}
