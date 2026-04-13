package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1305",
		Title:    "Avoid `$COMP_WORDS` — use `$words` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_WORDS` is a Bash completion variable containing the words on " +
			"the command line. Zsh completion uses `$words` array for the same purpose.",
		Check: checkZC1305,
	})
}

func checkZC1305(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$COMP_WORDS" && ident.Value != "COMP_WORDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1305",
		Message: "Avoid `$COMP_WORDS` in Zsh — use `$words` array instead. `COMP_WORDS` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
