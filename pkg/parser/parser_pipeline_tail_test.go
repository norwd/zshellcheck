// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParsePipelineTailAfterFor(t *testing.T) {
	parseSourceClean(t, "for f in *; do echo $f; done | wc -l\n")
}

func TestParsePipelineTailAfterIf(t *testing.T) {
	parseSourceClean(t, "if true; then echo yes; fi | tee log\n")
}

func TestParsePipelineTailAfterWhile(t *testing.T) {
	parseSourceClean(t, "while read -r l; do echo $l; done < f | sort\n")
}

func TestParsePipelineTailAfterCase(t *testing.T) {
	parseSourceClean(t, "case $x in a) echo a;; esac | tr a-z A-Z\n")
}

func TestParsePipelineTailLogicalOr(t *testing.T) {
	parseSourceClean(t, "for f in *; do echo $f; done || echo done\n")
}
