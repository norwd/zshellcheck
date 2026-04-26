// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

func tokensFor(t *testing.T, src string) []token.Token {
	t.Helper()
	l := New(src)
	var out []token.Token
	for {
		tok := l.NextToken()
		out = append(out, tok)
		if tok.Type == token.EOF {
			return out
		}
		if len(out) > 4096 {
			t.Fatalf("runaway lexer for %q", src)
		}
	}
}

func TestHeredocBasic(t *testing.T) {
	tokensFor(t, "cat <<EOF\nhello\nEOF\n")
}

func TestHeredocStripTabs(t *testing.T) {
	tokensFor(t, "cat <<-EOF\n\thello\n\tEOF\n")
}

func TestHeredocQuotedDelimiter(t *testing.T) {
	tokensFor(t, "cat <<\"END\"\n$x\nEND\n")
}

func TestHeredocSingleQuotedDelimiter(t *testing.T) {
	tokensFor(t, "cat <<'END'\nliteral\nEND\n")
}

func TestHeredocBackslashDelimiter(t *testing.T) {
	tokensFor(t, "cat <<\\EOF\nliteral\nEOF\n")
}

func TestHeredocUnterminated(t *testing.T) {
	tokensFor(t, "cat <<EOF\nthe end never comes\n")
}

func TestHeredocNoDelimiter(t *testing.T) {
	tokensFor(t, "cat <<\n")
}

func TestHeredocStripTabsNoBody(t *testing.T) {
	tokensFor(t, "cat <<-\n")
}
