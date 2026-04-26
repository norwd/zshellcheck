// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package lexer

import "testing"

func TestBackslashEscapeIllegalChar(t *testing.T) {
	tokensFor(t, "echo \\\x01\n")
}

func TestBackslashEscapeAtEOF(t *testing.T) {
	tokensFor(t, "echo \\")
}

func TestStringWithBraceParameterExpansion(t *testing.T) {
	tokensFor(t, "echo \"${var=\"default\"}\"\n")
}

func TestStringSingleQuoteLiteralBackslash(t *testing.T) {
	tokensFor(t, "echo 'no\\escape'\n")
}

func TestStringAnsiCQuoted(t *testing.T) {
	tokensFor(t, "echo $'tab\\there'\n")
}

func TestStringDoubleQuoteUnterminated(t *testing.T) {
	tokensFor(t, "echo \"never ends\n")
}

func TestStringSingleQuoteUnterminated(t *testing.T) {
	tokensFor(t, "echo 'never ends\n")
}

func TestPeekAtBeyondEnd(t *testing.T) {
	l := New("ab")
	if got := l.peekAt(99); got != 0 {
		t.Errorf("peekAt past end should return 0, got %d", got)
	}
	// peekAt(1) is the next byte after the cursor. New("ab") starts
	// pointing at 'a'; readPosition == 1 so peekAt(1) == 'b'.
	_ = l.peekAt(1)
}

func TestAngleBracketOperatorVariants(t *testing.T) {
	for _, src := range []string{
		"echo a >| b\n",
		"echo a >& b\n",
		"echo a 2>&1\n",
		"echo a <(cmd)\n",
		"echo a >(cmd)\n",
		"echo a <<<<EOF\nhello\nEOF\n",
	} {
		tokensFor(t, src)
	}
}
