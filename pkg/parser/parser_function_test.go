// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseFunctionLiteralKeywordForm(t *testing.T) {
	parseSourceClean(t, "function name() { echo hi; }\n")
}

func TestParseFunctionLiteralBareForm(t *testing.T) {
	parseSourceClean(t, "name() { echo hi; }\n")
}

func TestParseFunctionLiteralKeywordWithoutParens(t *testing.T) {
	parseSourceClean(t, "function greet { echo hi; }\n")
}

func TestParseFunctionLiteralBody(t *testing.T) {
	parseSourceClean(t, "f() { local x=1; echo $x; }\n")
}
