// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestFlagArgPositionMatch(t *testing.T) {
	cmd := &ast.SimpleCommand{
		Token: token.Token{Type: token.IDENT, Literal: "gpg", Line: 1, Column: 1},
		Name:  &ast.Identifier{Token: token.Token{Literal: "gpg"}, Value: "gpg"},
		Arguments: []ast.Expression{
			&ast.Identifier{Token: token.Token{Literal: "--delete-secret-keys", Line: 1, Column: 5}, Value: "--delete-secret-keys"},
		},
	}
	line, col := FlagArgPosition(cmd, map[string]bool{"--delete-secret-keys": true})
	if line != 1 || col != 5 {
		t.Errorf("expected (1,5), got (%d,%d)", line, col)
	}
}

func TestFlagArgPositionNoMatch(t *testing.T) {
	cmd := &ast.SimpleCommand{
		Token: token.Token{Type: token.IDENT, Literal: "ls", Line: 7, Column: 3},
		Name:  &ast.Identifier{Token: token.Token{Literal: "ls"}, Value: "ls"},
		Arguments: []ast.Expression{
			&ast.Identifier{Token: token.Token{Literal: "-l"}, Value: "-l"},
		},
	}
	line, col := FlagArgPosition(cmd, map[string]bool{"--secret": true})
	if line != 7 || col != 3 {
		t.Errorf("expected fallback to cmd token (7,3), got (%d,%d)", line, col)
	}
}

func TestFlagArgPositionEmptyArgs(t *testing.T) {
	cmd := &ast.SimpleCommand{
		Token: token.Token{Type: token.IDENT, Literal: "gpg", Line: 9, Column: 9},
		Name:  &ast.Identifier{Token: token.Token{Literal: "gpg"}, Value: "gpg"},
	}
	line, col := FlagArgPosition(cmd, map[string]bool{"--delete": true})
	if line != 9 || col != 9 {
		t.Errorf("expected fallback (9,9), got (%d,%d)", line, col)
	}
}

func TestIsIdentByteAllRanges(t *testing.T) {
	for _, b := range []byte{'a', 'z', 'A', 'Z', '0', '9', '_', '-'} {
		if !isIdentByte(b) {
			t.Errorf("expected %q to be ident byte", b)
		}
	}
	for _, b := range []byte{' ', '\t', '!', '#', '$', '/', '@'} {
		if isIdentByte(b) {
			t.Errorf("expected %q to NOT be ident byte", b)
		}
	}
}
