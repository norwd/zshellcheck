package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1318",
		Title:    "Avoid `$BASH_CMDS` — use `$commands` hash in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_CMDS` is a Bash associative array caching command lookups. " +
			"Zsh provides the `$commands` hash for the same purpose, mapping " +
			"command names to their full paths.",
		Check: checkZC1318,
	})
}

func checkZC1318(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_CMDS" && ident.Value != "BASH_CMDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1318",
		Message: "Avoid `$BASH_CMDS` in Zsh — use the `$commands` hash for command path lookups instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
