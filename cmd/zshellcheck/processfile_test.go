// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// dirty contains a backtick command sub which ZC1002 flags + auto-fixes.
const dirty = "#!/bin/zsh\nresult=`which git`\necho $result\n"

func TestProcessFileText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{})
}

func TestProcessFileJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "json", nil, fixOptions{})
}

func TestProcessFileSarif(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "sarif", nil, fixOptions{})
}

func TestProcessFileNonexistent(t *testing.T) {
	var out, errOut bytes.Buffer
	got := processFile("/nonexistent/path/zzz.zsh", &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{})
	if got != 0 {
		t.Errorf("expected 0 violations on read error, got %d", got)
	}
}

func TestProcessFileFixDryRun(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	stats := &fixStats{}
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{enabled: true, dryRun: true, stats: stats})
}

func TestProcessFileFixDiff(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{enabled: true, diff: true, dryRun: true})
}

func TestProcessFileFixApply(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	stats := &fixStats{}
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{enabled: true, stats: stats})
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) == dirty {
		t.Errorf("expected file rewrite, got unchanged contents")
	}
}

func TestProcessFileSeverityFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte(dirty), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	allowed := []katas.Severity{katas.SeverityError}
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", allowed, fixOptions{})
}
