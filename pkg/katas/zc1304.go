package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1304",
		Title:    "Avoid `$BASH_SUBSHELL` — use `$ZSH_SUBSHELL` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_SUBSHELL` tracks subshell nesting depth in Bash. " +
			"Zsh provides `$ZSH_SUBSHELL` as the native equivalent.",
		Check: checkZC1304,
	})
}

func checkZC1304(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_SUBSHELL" && ident.Value != "BASH_SUBSHELL" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1304",
		Message: "Avoid `$BASH_SUBSHELL` in Zsh — use `$ZSH_SUBSHELL` instead. `BASH_SUBSHELL` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}
