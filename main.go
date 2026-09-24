package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const obsidianRepo = "https://github.com/kepano/obsidian-skills"

var obsidianSkills = []string{"obsidian-bases", "obsidian-markdown"}

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "install" && os.Args[1] != "uninstall") {
		usage()
	}
	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	withObsidian := fs.Bool("obsidian", true, "also install "+fmt.Sprint(obsidianSkills)+" from kepano/obsidian-skills; -obsidian=false to skip")
	fs.Parse(os.Args[2:])
	if fs.NArg() != 0 {
		usage()
	}

	src, err := filepath.Abs("skills")
	if err != nil {
		fail(err)
	}
	if _, err := os.Stat(src); err != nil {
		fail(fmt.Errorf("run from the repo root: %w", err))
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fail(err)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		fail(err)
	}
	dst := filepath.Join(home, ".claude", "skills")
	vendor := filepath.Join(cache, "jaiminho-brain", "obsidian-skills")

	if os.Args[1] == "uninstall" {
		err = uninstall(dst, src, filepath.Join(vendor, "skills"))
	} else if err = link(src, dst, nil); err == nil && *withObsidian {
		if err = fetch(obsidianRepo, "HEAD", vendor); err == nil {
			err = link(filepath.Join(vendor, "skills"), dst, obsidianSkills)
		}
	}
	if err != nil {
		fail(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run . install [-obsidian=false] | uninstall")
	os.Exit(2)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func fetch(repo, rev, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"fetch", "-q", "--depth", "1", repo, rev},
		{"checkout", "-q", "--force", "FETCH_HEAD"},
		{"log", "-1", "--format=%h %cs %s"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git %v: %w", args, err)
		}
	}
	return nil
}

func link(src, dst string, only []string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	names := only
	if names == nil {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}
	}
	for _, name := range names {
		target := filepath.Join(src, name)
		if _, err := os.Stat(filepath.Join(target, "SKILL.md")); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		l := filepath.Join(dst, name)
		if fi, err := os.Lstat(l); err == nil {
			if fi.Mode()&os.ModeSymlink == 0 {
				fmt.Fprintf(os.Stderr, "skip %s: %s exists and is not a symlink\n", name, l)
				continue
			}
			if err := os.Remove(l); err != nil {
				return err
			}
		}
		if err := os.Symlink(target, l); err != nil {
			return err
		}
		fmt.Println("linked", name)
	}
	return nil
}

func uninstall(dst string, srcs ...string) error {
	entries, err := os.ReadDir(dst)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, e := range entries {
		l := filepath.Join(dst, e.Name())
		target, err := os.Readlink(l)
		if err != nil {
			continue
		}
		for _, src := range srcs {
			if filepath.Dir(target) == src {
				if err := os.Remove(l); err != nil {
					return err
				}
				fmt.Println("unlinked", e.Name())
				break
			}
		}
	}
	return nil
}
