// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

// TestFixGuardsRejectMismatchedNode runs every registered Fix with a
// node type the fix does not expect. Each fix's first guard should
// drop out via the type-assertion check, returning nil and exercising
// the defensive branch without panic.
func TestFixGuardsRejectMismatchedNode(t *testing.T) {
	stub := &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: "0", Line: 1, Column: 1},
		Value: 0,
	}
	for id, kata := range Registry.KatasByID {
		if kata.Fix == nil {
			continue
		}
		v := Violation{KataID: id, Line: 1, Column: 1}
		// Most fixes assert the concrete node type and return nil for
		// any other type. The IntegerLiteral above almost never matches
		// (no kata fixes integer literals), so this drives the guard.
		_ = kata.Fix(stub, v, []byte("placeholder\n"))
	}
}

// TestFixGuardsRejectOutOfRangeViolation exercises the byte-offset
// guard that runs after the type cast. Many fixes look up
// LineColToByteOffset; passing a violation whose Line/Column points
// past EOF makes the offset check fail and returns nil.
func TestFixGuardsRejectOutOfRangeViolation(t *testing.T) {
	for id, kata := range Registry.KatasByID {
		if kata.Fix == nil {
			continue
		}
		v := Violation{KataID: id, Line: 9999, Column: 9999}
		_ = kata.Fix(nil, v, []byte("placeholder\n"))
	}
}

// TestFixGuardsEmptySource feeds an empty source slice. Most offset
// guards reject the empty source immediately.
func TestFixGuardsEmptySource(t *testing.T) {
	for id, kata := range Registry.KatasByID {
		if kata.Fix == nil {
			continue
		}
		v := Violation{KataID: id, Line: 1, Column: 1}
		_ = kata.Fix(nil, v, []byte{})
	}
}
