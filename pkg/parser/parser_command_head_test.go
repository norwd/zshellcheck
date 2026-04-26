// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseCommandHeadDollarParen(t *testing.T) {
	parseSourceClean(t, "$(echo hello) world\n")
}

func TestParseCommandHeadBacktick(t *testing.T) {
	parseSourceClean(t, "`echo hi` world\n")
}

func TestParseCommandHeadVariable(t *testing.T) {
	parseSourceClean(t, "$cmd arg1 arg2\n")
}

func TestParseCommandHeadDollarBrace(t *testing.T) {
	parseSourceClean(t, "${cmd:-default} arg\n")
}

func TestParseCommandWithLeadingParenArg(t *testing.T) {
	parseSourceClean(t, "echo (foo bar)\n")
}

func TestParseCommandFunctionDefinition(t *testing.T) {
	parseSourceClean(t, "myfunc() { echo hi; }\n")
}

func TestParseCommandFunctionDefinitionMultiline(t *testing.T) {
	parseSourceClean(t, "myfunc() {\n  echo hi\n  echo bye\n}\n")
}

func TestParseCommandPipelineMulti(t *testing.T) {
	parseSourceClean(t, "ls | sort | uniq | wc -l\n")
}

func TestParseCommandPipelineNegated(t *testing.T) {
	parseSourceClean(t, "! grep foo file\n")
}

func TestParseCommandLogicalChain(t *testing.T) {
	parseSourceClean(t, "true && echo yes || echo no\n")
}

func TestParseCommandWithBraceArg(t *testing.T) {
	parseSourceClean(t, "echo {a,b,c}.zsh\n")
}

func TestParseCommandConcatenatedString(t *testing.T) {
	parseSourceClean(t, "echo \"a\"\"b\"\"c\"\n")
}
