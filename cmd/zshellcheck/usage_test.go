// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"strings"
	"testing"
)

func TestPaletteDisabled(t *testing.T) {
	p := palette{enabled: false}
	if got := p.bold("x"); got != "x" {
		t.Errorf("disabled bold returns wrapped: %q", got)
	}
	if got := p.dim("x"); got != "x" {
		t.Errorf("disabled dim returns wrapped: %q", got)
	}
	if got := p.section("x"); got != "x" {
		t.Errorf("disabled section returns wrapped: %q", got)
	}
}

func TestPaletteEnabled(t *testing.T) {
	p := palette{enabled: true}
	out := p.bold("x")
	if !strings.Contains(out, "\x1b[1m") || !strings.Contains(out, "\x1b[0m") {
		t.Errorf("enabled bold missing ANSI: %q", out)
	}
}
