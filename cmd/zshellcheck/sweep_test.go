// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// kitchen sink fixture: each line targets a different kata pattern.
// processFile runs Check + Fix on every node so the fixer paths
// matter even when no assertions are made.
const sweepFixture = `#!/usr/bin/env zsh
result=` + "`which git`" + `
echo $arr[1]
target=$1
echo -E "Cleaning $target"
rm -rf $target
joined=$(seq -s, 1 5)
for f in *; do echo $f; done | wc -l
if [ -f config ]; then echo yes; fi
[[ -z $foo ]] && echo empty
typeset -a items=(a b c)
local x=$(echo nested)
function greet() { echo hello; }
case $x in a) echo a;; esac
arr[(R)x]=1
echo "${arr[@]}"
n=$(( 1 + 2 ))
trap 'echo bye' EXIT
read -r line < input
print -r -- "${(j:,:)items}"
`

func TestSweepFixApply(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sweep.zsh")
	if err := os.WriteFile(path, []byte(sweepFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	stats := &fixStats{}
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{enabled: true, maxPasses: 5, stats: stats})
}

func TestSweepFixDiff(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sweep.zsh")
	if err := os.WriteFile(path, []byte(sweepFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{enabled: true, diff: true, dryRun: true, maxPasses: 5})
}

func TestSweepFormats(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sweep.zsh")
	if err := os.WriteFile(path, []byte(sweepFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, fmt := range []string{"text", "json", "sarif"} {
		var out, errOut bytes.Buffer
		processFile(path, &out, &errOut, config.DefaultConfig(), katas.Registry, fmt, nil, fixOptions{})
	}
}
