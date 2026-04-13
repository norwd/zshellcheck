package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1307",
		Title:    "Avoid `$DIRSTACK` — use `$dirstack` (lowercase) in Zsh",
		Severity: SeverityWarning,
		Description: "`$DIRSTACK` is the Bash form of the directory stack array. " +
			"Zsh uses `$dirstack` (lowercase) for the same purpose.",
		Check: checkZC1307,
	})
}

func checkZC1307(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$DIRSTACK" && ident.Value != "DIRSTACK" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1307",
		Message: "Avoid `$DIRSTACK` in Zsh — use `$dirstack` (lowercase) instead. The uppercase form is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
