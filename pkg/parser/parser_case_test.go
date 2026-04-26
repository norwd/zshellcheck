// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseCaseStatementBasic(t *testing.T) {
	parseSourceClean(t, "case $x in a) echo a;; b) echo b;; esac\n")
}

func TestParseCaseStatementMultiplePatterns(t *testing.T) {
	parseSourceClean(t, "case $x in a|b|c) echo abc;; *) echo other;; esac\n")
}

func TestParseCaseStatementParenLabel(t *testing.T) {
	parseSourceClean(t, "case $x in (a) echo a;; (b) echo b;; esac\n")
}

func TestParseCaseStatementGlobAlternation(t *testing.T) {
	parseSourceClean(t, "case $x in (darwin|freebsd)*) echo bsd;; esac\n")
}

func TestParseCaseStatementNested(t *testing.T) {
	parseSourceClean(t, "case x in a) case y in 1) echo nested;; esac;; esac\n")
}

func TestParseCaseStatementEmptyClauses(t *testing.T) {
	parseSourceClean(t, "case $x in a) ;; b) ;; esac\n")
}

func TestParseAnonymousFunction(t *testing.T) {
	parseSourceClean(t, "() { echo anon; }\n")
}

func TestParseShebang(t *testing.T) {
	parseSourceClean(t, "#!/usr/bin/env zsh\necho ok\n")
}

func TestParseDoubleBracketTest(t *testing.T) {
	parseSourceClean(t, "[[ -f file && -r file ]]\n")
}

func TestParseDoubleBracketRegex(t *testing.T) {
	parseSourceClean(t, "[[ $x =~ ^[a-z]+$ ]]\n")
}

func TestParseProcessSubstitution(t *testing.T) {
	parseSourceClean(t, "diff <(sort a) <(sort b)\n")
}

func TestParseSelectStatement(t *testing.T) {
	parseSourceClean(t, "select opt in a b c; do echo $opt; break; done\n")
}
