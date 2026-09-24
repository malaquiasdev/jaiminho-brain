package main

import (
	"os"
	"path/filepath"
	"testing"
)

func mkskills(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.MkdirAll(filepath.Join(root, n), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, n, "SKILL.md"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLinkUninstall(t *testing.T) {
	src := filepath.Join(t.TempDir(), "skills")
	vendor := filepath.Join(t.TempDir(), "skills")
	dst := t.TempDir()
	mkskills(t, src, "brain-log", "mine")
	mkskills(t, vendor, "obsidian-bases", "json-canvas")
	if err := os.Mkdir(filepath.Join(dst, "mine"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/elsewhere/other", filepath.Join(dst, "other")); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := link(src, dst, nil); err != nil {
			t.Fatal(err)
		}
		if err := link(vendor, dst, []string{"obsidian-bases"}); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := os.Readlink(filepath.Join(dst, "brain-log")); got != filepath.Join(src, "brain-log") {
		t.Fatalf("brain-log -> %q", got)
	}
	if fi, _ := os.Lstat(filepath.Join(dst, "mine")); fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal("real dir was replaced")
	}
	if _, err := os.Lstat(filepath.Join(dst, "json-canvas")); !os.IsNotExist(err) {
		t.Fatal("json-canvas linked but was not selected")
	}
	if err := link(vendor, dst, []string{"missing"}); err == nil {
		t.Fatal("linking a missing skill should fail")
	}

	if err := uninstall(dst, src, vendor); err != nil {
		t.Fatal(err)
	}
	for _, gone := range []string{"brain-log", "obsidian-bases"} {
		if _, err := os.Lstat(filepath.Join(dst, gone)); !os.IsNotExist(err) {
			t.Fatalf("%s still linked", gone)
		}
	}
	for _, keep := range []string{"mine", "other"} {
		if _, err := os.Lstat(filepath.Join(dst, keep)); err != nil {
			t.Fatalf("%s removed: %v", keep, err)
		}
	}
}
