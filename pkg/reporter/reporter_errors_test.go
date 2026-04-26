// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"errors"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// limitWriter fails after writing n bytes — drives Report's
// io.Writer error branches.
type limitWriter struct {
	remaining int
}

func (w *limitWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errors.New("writer full")
	}
	if len(p) > w.remaining {
		written := w.remaining
		w.remaining = 0
		return written, errors.New("writer full mid-write")
	}
	w.remaining -= len(p)
	return len(p), nil
}

func TestTextReporter_LocationWriteError(t *testing.T) {
	w := &limitWriter{remaining: 0}
	r := NewTextReporter(w, "f.zsh", "echo hello\n", config.Config{NoColor: true})
	err := r.Report([]katas.Violation{{KataID: "ZC1", Line: 1, Column: 1, Level: katas.SeverityError}})
	if err == nil {
		t.Error("expected error from writer")
	}
}

func TestTextReporter_SeverityHeaderError(t *testing.T) {
	w := &limitWriter{remaining: 12}
	r := NewTextReporter(w, "f.zsh", "echo hello\n", config.Config{NoColor: true})
	err := r.Report([]katas.Violation{{KataID: "ZC1", Line: 1, Column: 1, Level: katas.SeverityError, Message: "msg"}})
	if err == nil {
		t.Error("expected error from writer")
	}
}

func TestTextReporter_SnippetWriteError(t *testing.T) {
	w := &limitWriter{remaining: 36}
	r := NewTextReporter(w, "f.zsh", "echo hello\n", config.Config{NoColor: true})
	err := r.Report([]katas.Violation{{KataID: "ZC1", Line: 1, Column: 1, Level: katas.SeverityError, Message: "msg"}})
	if err == nil {
		t.Error("expected error from writer")
	}
}

func TestTextReporter_NegativeColumnPadding(t *testing.T) {
	w := &limitWriter{remaining: 4096}
	r := NewTextReporter(w, "f.zsh", "echo hello\n", config.Config{NoColor: true})
	err := r.Report([]katas.Violation{{KataID: "ZC1", Line: 1, Column: -5, Level: katas.SeverityInfo, Message: "neg col"}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTextReporter_AllSeverityColors(t *testing.T) {
	w := &limitWriter{remaining: 1 << 16}
	r := NewTextReporter(w, "f.zsh", "echo hello\n", config.Config{NoColor: false})
	violations := []katas.Violation{
		{KataID: "E", Line: 1, Column: 1, Level: katas.SeverityError, Message: "e"},
		{KataID: "W", Line: 1, Column: 1, Level: katas.SeverityWarning, Message: "w"},
		{KataID: "I", Line: 1, Column: 1, Level: katas.SeverityInfo, Message: "i"},
		{KataID: "S", Line: 1, Column: 1, Level: katas.SeverityStyle, Message: "s"},
	}
	if err := r.Report(violations); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
