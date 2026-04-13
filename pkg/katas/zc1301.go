package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1301",
		Title:    "Avoid `$PIPESTATUS` — use `$pipestatus` (lowercase) in Zsh",
		Severity: SeverityWarning,
		Description: "`$PIPESTATUS` is a Bash array containing exit statuses from the last " +
			"pipeline. Zsh uses `$pipestatus` (lowercase) for the same purpose. " +
			"The uppercase form is undefined in Zsh.",
		Check: checkZC1301,
	})
}

func checkZC1301(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$PIPESTATUS" && ident.Value != "PIPESTATUS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1301",
		Message: "Avoid `$PIPESTATUS` in Zsh — use `$pipestatus` (lowercase) instead. The uppercase form is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
