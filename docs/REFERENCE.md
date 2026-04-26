# Reference

Glossary, governance, and a comparison with related tooling.

## Contents

- [Governance](#governance)
- [Comparison vs ShellCheck](#comparison-vs-shellcheck)
- [Glossary](#glossary)

---

## Governance

ZShellCheck is led by **@afadesigns** (Andreas Fahl), the lead maintainer and sole committer.
The model below describes how community contribution scales as the project grows.

### Roles

- **Lead maintainer.**
  Final say on roadmap, release cadence, and architectural direction.
  Held by @afadesigns.
- **Maintainers.**
  Review and merge PRs, triage issues.
  CODEOWNERS routes every review to @afadesigns.
- **Contributors.**
  Anyone who opens an issue or PR.

### Decision making

Non-trivial direction is discussed on GitHub issues or PRs.
Consensus is preferred; the lead maintainer has the casting vote when consensus stalls.
Breaking changes and major features require a design-discussion issue before implementation.

### Review requirements

- Every PR to `main` needs one approving review from the CODEOWNERS set.
- Commits must be GPG-signed.
- Required status checks — `test`, `security`, `sbom` — must pass.
- No force-push, no branch deletion, no unsigned merges.

---

## Comparison vs ShellCheck

| Feature | ShellCheck | ZShellCheck |
| :--- | :--- | :--- |
| Focus | `sh` and `bash` (POSIX) | `zsh` (native) |
| Language | Haskell | Go |
| Philosophy | Portability | Zsh power |
| Checks | ~500 | 1000 |
| Output | Text, JSON, GCC, TTY | Text, JSON, SARIF |
| Severity | error, warning, info, style | error, warning, info, style |
| Auto-fix | Partial | First-class — `-fix`, `-diff`, `-dry-run`. The fix-enabled count appears in [KATAS.md](../KATAS.md). |

Use ZShellCheck for `.zshrc`, Zsh plugins, and scripts that lean on Zsh-specific features such as `[[`, parameter modifiers, glob qualifiers, and associative arrays.

---

## Glossary

- **Kata.**
  A specific check or rule, identified by a `ZC` prefix and a four-digit number — for example `ZC1001`.
- **AST.**
  Abstract syntax tree — the structured representation of parsed Zsh source.
- **Lexer.**
  The tokenizer; converts source text into a stream of tokens.
- **Walker.**
  An AST traverser; visits each node so katas can apply.
- **Registry.**
  The central store of katas indexed by node type.
- **Severity.**
  The impact level of a kata violation: `error`, `warning`, `info`, or `style`.
- **SARIF.**
  Static Analysis Results Interchange Format; the JSON schema used for GitHub Code Scanning ingestion.
