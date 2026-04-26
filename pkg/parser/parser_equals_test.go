// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

// parseSourceClean parses src and fails the test on any parser error.
func parseSourceClean(t *testing.T, src string) *Parser {
	t.Helper()
	p := New(lexer.New(src))
	prog := p.ParseProgram()
	if prog == nil {
		t.Fatalf("ParseProgram returned nil for %q", src)
	}
	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected parser errors for %q: %v", src, errs)
	}
	return p
}

func TestParseEqualsForm(t *testing.T) {
	parseSourceClean(t, "=ls -la\n")
}

func TestParseEqualsFormSpaceTerminates(t *testing.T) {
	// `= ls` should not absorb `ls` because peek has preceding space.
	p := New(lexer.New("= ls\n"))
	_ = p.ParseProgram()
}
