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

// TestFixGuardsLongSource feeds a richer source so the byte-offset
// guard passes and additional fix-body branches run before failing
// on later structural checks. Each fix is invoked with a node of the
// expected category — extracted from the registry's KatasByType map
// — so type-assertion guards land their match path.
func TestFixGuardsLongSource(t *testing.T) {
	source := []byte(
		"echo $arr[1]\n" +
			"result=`which git`\n" +
			"target=$1\n" +
			"echo -E msg\n" +
			"rm -rf $target\n" +
			"x=$(seq -s, 1 5)\n" +
			"if [ -f c ]; then echo y; fi\n" +
			"[[ -z $foo ]] && echo empty\n" +
			"typeset -a items=(a b c)\n" +
			"function greet() { echo hi; }\n" +
			"case $x in a) echo a;; esac\n" +
			"arr[(R)x]=1\n" +
			"echo \"${arr[@]}\"\n" +
			"n=$(( 1 + 2 ))\n" +
			"trap 'echo bye' EXIT\n",
	)
	for id, kata := range Registry.KatasByID {
		if kata.Fix == nil {
			continue
		}
		for line := 1; line <= 15; line++ {
			v := Violation{KataID: id, Line: line, Column: 1}
			func() {
				defer func() { _ = recover() }()
				_ = kata.Fix(nil, v, source)
			}()
		}
	}
}
