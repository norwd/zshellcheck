// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"regexp"
	"sort"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// Metamorphic property: semantic-preserving rewrites of a Zsh script must not
// change which katas fire. The parser-corpus and violation-corpus sweeps prove
// the linter is stable on a fixed corpus; this proves it is stable under
// reformatting of the same code. A kata whose finding appears or disappears
// when blank lines or comments are added around inert code is position- or
// whitespace-sensitive, which is a structural false-positive risk. Comparing
// the kata ID multiset, not line numbers, isolates which rules fired from where.

// kataIDMultiset returns the sorted kata IDs, duplicates included, for a source.
func kataIDMultiset(code string) []string {
	vs := testutil.CheckAll(code)
	ids := make([]string, 0, len(vs))
	for _, v := range vs {
		ids = append(ids, v.KataID)
	}
	sort.Strings(ids)
	return ids
}

// metamorphicSamples are inert-context Zsh snippets that each fire one or more
// katas. None carries a shebang or a here-document, so wrapping them in leading
// or trailing blank lines and comments is a guaranteed semantic no-op.
var metamorphicSamples = []string{
	`echo -E "Cleaning $target"`,
	`rm -rf $target`,
	`echo $my_array[1]`,
	"result=`date`",
	`for i in $(seq 1 10); do echo $i; done`,
	`if [ "$a" = "$b" ]; then echo equal; fi`,
	`typeset -A map; map[key]=value`,
	`cat file | grep needle`,
	`[[ -z $var ]] && echo empty`,
	`x=$((1 + 2)); echo $x`,
	`ls | wc -l`,
	`export PATH=$PATH:/usr/local/bin`,
}

// preservingTransforms each return code that is semantically identical to the
// input for analysis purposes: only blank lines and whole-line comments are
// added before or after the snippet, never inside it.
var preservingTransforms = []struct {
	name string
	fn   func(string) string
}{
	{"identity", func(s string) string { return s }},
	{"prepend-blanks", func(s string) string { return "\n\n\n" + s }},
	{"append-blanks", func(s string) string { return s + "\n\n\n" }},
	{"prepend-comment", func(s string) string { return "# leading comment\n" + s }},
	{"append-comment", func(s string) string { return s + "\n# trailing comment\n" }},
	{"wrap-comments", func(s string) string { return "# top\n\n" + s + "\n\n# bottom\n" }},
	{"trailing-newline", func(s string) string { return s + "\n" }},
}

func TestMetamorphicFormatInvariance(t *testing.T) {
	for _, sample := range metamorphicSamples {
		want := kataIDMultiset(sample)
		for _, tr := range preservingTransforms {
			got := kataIDMultiset(tr.fn(sample))
			if !equalStringSlices(want, got) {
				t.Errorf("transform %q changed the kata findings\n  sample:   %q\n  baseline: %v\n  got:      %v",
					tr.name, sample, want, got)
			}
		}
	}
}

// TestMetamorphicRenameInvariance asserts that consistently renaming a plain
// user variable does not change which katas fire. The variable names are
// neutral, none is a name a kata special-cases (no PATH, RANDOM, password,
// secret, and the like), so the rename is semantically inert and the finding
// set must be identical. A change here means a kata keys on a specific spelling
// where it should key on structure, which is a false-positive risk.
func TestMetamorphicRenameInvariance(t *testing.T) {
	cases := []struct{ code, varName string }{
		{`echo $myarr[1]`, "myarr"},
		{`rm -rf $dir`, "dir"},
		{`for item in $(seq 1 3); do echo $item; done`, "item"},
		{`lhs=1; rhs=2; [ "$lhs" = "$rhs" ] && echo eq`, "lhs"},
		{`val=$((1 + 2)); echo $val`, "val"},
	}
	for _, c := range cases {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(c.varName) + `\b`)
		renamed := re.ReplaceAllString(c.code, c.varName+"_renamed")
		want := kataIDMultiset(c.code)
		got := kataIDMultiset(renamed)
		if !equalStringSlices(want, got) {
			t.Errorf("renaming %q changed the kata findings\n  code:     %q\n  renamed:  %q\n  baseline: %v\n  got:      %v",
				c.varName, c.code, renamed, want, got)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
