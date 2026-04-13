package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1308",
		Title:    "Avoid `$COMP_LINE` — use `$BUFFER` in Zsh completion",
		Severity: SeverityWarning,
		Description: "`$COMP_LINE` is a Bash completion variable containing the full command " +
			"line. Zsh completion uses `$BUFFER` for the current command line content.",
		Check: checkZC1308,
	})
}

func checkZC1308(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$COMP_LINE" && ident.Value != "COMP_LINE" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1308",
		Message: "Avoid `$COMP_LINE` in Zsh — use `$BUFFER` instead. `COMP_LINE` is Bash completion-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
