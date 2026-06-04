// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

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

// A double-quoted string whose `${…}` body contains an escaped `\${`
// must not let the lexer count the following `{` as opening a nested
// expansion. zsh closes the expansion at the next unescaped `}` and a
// lone `{` is literal, so the closing `"` must terminate the string.
// Before the fix the inflated brace depth swallowed the quote and the
// rest of the file cascaded into one giant string (powerlevel10k
// p10k.zsh line 6786). Issue #1377.
func TestStringEscapedDollarBraceDoesNotSwallow(t *testing.T) {
	for _, src := range []string{
		"echo \"${(%):-a\\${b}\"\nnext=1\n",
		"echo \"${x:-a{b}\"\nnext=2\n",
		"echo \"${(%):-  %3F'(( ! \\${+functions[p10k]\\} )) || x'%f >>! $V}\"\nnext=3\n",
	} {
		toks := tokensFor(t, src)
		var sawString, sawNext bool
		for _, tk := range toks {
			if tk.Type == token.STRING {
				sawString = true
			}
			if tk.Type == token.IDENT && (tk.Literal == "next") {
				sawNext = true
			}
		}
		if !sawString {
			t.Errorf("no STRING token for %q", src)
		}
		// The code after the string must survive — if the closing quote
		// were swallowed, `next=N` would be lexed as string content.
		if !sawNext {
			t.Errorf("trailing assignment swallowed (quote not closed) for %q", src)
		}
	}
}

// The nested-quote skip inside `${…}` covers a `'…'` span, a `"…"`
// span with a `\"` escape, and an unterminated span that runs to EOF.
// A brace hidden in any of these must not unbalance the expansion.
func TestStringNestedQuoteSpansInsideDollarBrace(t *testing.T) {
	for _, src := range []string{
		"echo \"${x:-'a}b'}\"\nnext=1\n",       // single-quoted span hides `}`
		"echo \"${x:-\"a\\\"}b\"}\"\nnext=2\n", // double-quoted span with \" and `}`
		"echo \"${x:-\"unterminated}",          // nested span runs to EOF
	} {
		// Must not panic or loop; tokensFor caps runaway lexing.
		tokensFor(t, src)
	}
}

// A nested `${…${…}…}` must still balance: the inner `${` opener is
// counted and its `}` decrements, so the outer expansion closes
// correctly. Guards against over-correcting the #1377 fix.
func TestStringNestedDollarBraceStillBalances(t *testing.T) {
	toks := tokensFor(t, "echo \"${x:-${y:-z}}\"\nafter=1\n")
	var sawAfter bool
	for _, tk := range toks {
		if tk.Type == token.IDENT && tk.Literal == "after" {
			sawAfter = true
		}
	}
	if !sawAfter {
		t.Error("nested ${…${…}…} mis-closed: trailing code swallowed")
	}
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
