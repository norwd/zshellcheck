// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseArithmeticTernary(t *testing.T) {
	parseSourceClean(t, "(( x = a > b ? a : b ))\n")
}

func TestParseArithmeticFloatLiteral(t *testing.T) {
	parseSourceClean(t, "(( x = 1.0 + 2.5 ))\n")
}

func TestParseArithmeticTrailingDot(t *testing.T) {
	parseSourceClean(t, "(( x = 1000. + 1 ))\n")
}

func TestParseArithmeticDollarParen(t *testing.T) {
	parseSourceClean(t, "x=$(( 1 + 2 ))\n")
}

func TestParseArithmeticIncrement(t *testing.T) {
	parseSourceClean(t, "(( i++ ))\n")
}
