// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseRedirectionStdoutTruncate(t *testing.T) {
	parseSourceClean(t, "echo hi > /dev/null\n")
}

func TestParseRedirectionStderrToStdout(t *testing.T) {
	parseSourceClean(t, "echo hi 2>&1\n")
}

func TestParseRedirectionAppend(t *testing.T) {
	parseSourceClean(t, "echo hi >> log\n")
}

func TestParseRedirectionStdinFromFile(t *testing.T) {
	parseSourceClean(t, "wc -l < input\n")
}

func TestParseRedirectionHerestring(t *testing.T) {
	parseSourceClean(t, "cat <<< hello\n")
}
