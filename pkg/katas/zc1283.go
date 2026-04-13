package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1283",
		Title:    "Use `setopt` instead of `set -o` for Zsh options",
		Severity: SeverityStyle,
		Description: "Zsh provides `setopt` and `unsetopt` as native builtins for managing shell " +
			"options. Using `set -o` / `set +o` is a POSIX compatibility form that is less " +
			"idiomatic in Zsh scripts.",
		Check: checkZC1283,
	})
}

func checkZC1283(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-o" {
			return []Violation{{
				KataID:  "ZC1283",
				Message: "Use `setopt` instead of `set -o` in Zsh scripts. `setopt` is the native Zsh idiom.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}
