<div align="center">

<img src="docs/assets/banner.png" alt="zshellcheck" width="100%" />

### The quiet linter for a quiet shell.

Static analysis and auto-fix for the setopts, hooks, and globs Bash never learned.

[![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/afadesigns/zshellcheck?color=blue)](https://github.com/afadesigns/zshellcheck/releases/latest)
[![Marketplace](https://img.shields.io/badge/Marketplace-ZshellCheck%20v1-2ea44f?logo=githubactions&logoColor=white)](https://github.com/marketplace/actions/zshellcheck-v1)
[![Auto-fix](https://img.shields.io/badge/auto--fix-67%20katas-2ea44f)](KATAS.md)
[![Go Report](https://goreportcard.com/badge/github.com/afadesigns/zshellcheck)](https://goreportcard.com/report/github.com/afadesigns/zshellcheck)
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg)](https://codecov.io/gh/afadesigns/zshellcheck)
[![Scorecard](https://api.securityscorecards.dev/projects/github.com/afadesigns/zshellcheck/badge)](https://securityscorecards.dev/viewer/?uri=github.com/afadesigns/zshellcheck)
[![SLSA](https://img.shields.io/badge/SLSA-Level%203-brightgreen)](https://slsa.dev)

</div>

---

## See it in action

<p align="center">
  <img src="docs/assets/demo.gif" alt="zshellcheck demo" width="100%" />
</p>

## Install

**The binary is the same wherever it lands.** Three ways to put it there:

```bash
# macOS, Linux, WSL
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

```powershell
# Windows
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 | iex
```

```bash
# Anywhere Go is installed
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

`--uninstall` reverses any of them.

**On every tag:** native `.deb`, `.rpm`, `.apk`, and a multi-arch container at `ghcr.io/afadesigns/zshellcheck`. The *full ledger* — pinning, cosign verification, distro one-liners — sits in [INSTALL.md](INSTALL.md).

## Run

```bash
zshellcheck path/to/script.zsh
zshellcheck -severity warning -format sarif ./scripts > zshellcheck.sarif
zshellcheck -diff path/to/script.zsh    # preview the auto-fix
zshellcheck -fix  path/to/script.zsh    # apply it
```

### GitHub Actions

```yaml
- uses: afadesigns/zshellcheck@v1.0.13
  with:
    args: -format sarif -severity warning ./scripts
```

### Pre-commit

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.13
    hooks:
      - id: zshellcheck
```


## Tested integrations

ZShellCheck is verified against the script corpora of these popular Zsh ecosystems.
Every release runs a parse + lint sweep over each — no panics, no crashes, deterministic output.

| Frameworks                                                                | Plugin / theme managers                                                            | Plugins + utilities                                                                                                                                                                                          | Themes / prompts                                                                                                          |
| :------------------------------------------------------------------------ | :--------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------ |
| oh-my-zsh · prezto · zimfw · antidote · zinit · zephyr · zsh-utils        | fzf · fzf-tab · fast-syntax-highlighting                                           | zsh-autosuggestions · zsh-syntax-highlighting · zsh-history-substring-search · zsh-vi-mode · zsh-autocomplete · zsh-completions                                                                                | powerlevel10k · spaceship-prompt · starship                                                                               |

## Documentation

| Doc | What's inside |
| --- | --- |
| [USER_GUIDE.md](docs/USER_GUIDE.md) | CLI reference, configuration, inline directives, integrations, FAQ |
| [DEVELOPER.md](docs/DEVELOPER.md) | Architecture, AST reference, kata authoring, release process |
| [REFERENCE.md](docs/REFERENCE.md) | Governance, glossary, ShellCheck comparison table |
| [KATAS.md](KATAS.md) | Every kata with description and severity |
| [CHANGELOG.md](CHANGELOG.md) | Per-release history |
| [SECURITY.md](SECURITY.md) | Vulnerability disclosure |
| [CONTRIBUTING.md](CONTRIBUTING.md) | PR workflow, local checks, conventions |
| [ROADMAP.md](ROADMAP.md) | LSP, auto-fixer, plugins |

## Contributing

PRs welcome.
Start with [CONTRIBUTING.md](CONTRIBUTING.md).
Questions and ideas: [discussions](https://github.com/afadesigns/zshellcheck/discussions).
Bugs: [issues](https://github.com/afadesigns/zshellcheck/issues).

## License

Distributed under the MIT License.
See [LICENSE](LICENSE).

## Credits

Andreas Fahl (**@afadesigns**) — author and maintainer.
Inspired by [ShellCheck](https://www.shellcheck.net/); ZShellCheck is an independent Go implementation focused on Zsh-specific semantics.

<div align="center">
  <a href="https://github.com/afadesigns/zshellcheck/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=afadesigns/zshellcheck" alt="Contributors" />
  </a>
</div>
