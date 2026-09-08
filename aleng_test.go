package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFixTerminfoEnv(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "78"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "78", "xterm-ghostty"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TERMINFO", dir)

	t.Setenv("TERM", "xterm-ghostty")
	fixTerminfoEnv()
	if os.Getenv("TERMINFO") != dir {
		t.Fatal("TERMINFO 에 있는 TERM 인데 TERMINFO 가 제거됨")
	}

	t.Setenv("TERM", "tmux-256color")
	fixTerminfoEnv()
	if os.Getenv("TERMINFO") != "" {
		t.Fatal("TERMINFO 에 없는 TERM 인데 TERMINFO 가 유지됨")
	}
}
