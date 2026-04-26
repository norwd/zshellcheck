// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"os"
	"testing"
)

func TestNewPaletteNonFileWriter(t *testing.T) {
	old := os.Getenv("NO_COLOR")
	defer func() { _ = os.Setenv("NO_COLOR", old) }()
	_ = os.Unsetenv("NO_COLOR")
	p := newPalette(&bytes.Buffer{})
	if p.enabled {
		t.Error("non-*os.File writer should produce disabled palette")
	}
}

func TestNewPaletteNoColorEnv(t *testing.T) {
	old := os.Getenv("NO_COLOR")
	defer func() { _ = os.Setenv("NO_COLOR", old) }()
	_ = os.Setenv("NO_COLOR", "1")
	p := newPalette(os.Stdout)
	if p.enabled {
		t.Error("NO_COLOR set should produce disabled palette")
	}
}

func TestNewPaletteRegularFile(t *testing.T) {
	old := os.Getenv("NO_COLOR")
	defer func() { _ = os.Setenv("NO_COLOR", old) }()
	_ = os.Unsetenv("NO_COLOR")
	dir := t.TempDir()
	f, err := os.Create(dir + "/log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p := newPalette(f)
	if p.enabled {
		t.Error("regular file (no char device) should produce disabled palette")
	}
}
