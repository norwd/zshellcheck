# Roadmap

ZShellCheck is an evolving static analysis tool for Zsh.
The goal is comprehensive, fast, and reliable tooling for the Zsh ecosystem.

## Milestones

### Version 0.0.x — foundation (complete)

- [x] Establish core architecture: lexer, parser, AST, walker.
- [x] Implement the linting framework.
- [x] Set up the CI/CD pipeline on GitHub Actions.
- [x] Create the initial kata set (ZC1000+).
- [x] Publish a documentation suite.
- [x] Ship the installer script with binary fallback.

### Version 0.1.x — growing the rule set (complete)

- [x] Severity levels (Error, Warning, Info, Style) on every kata.
- [x] `--severity` flag for filtering violations.
- [x] `--no-color` flag and verbose output mode.
- [x] SARIF output for GitHub Security integration.
- [x] Source-code context in violation output.
- [x] Installer with auto-sudo, CI detection, man pages, and completions.
- [x] 166 katas (ZC1001–ZC1169).
- [x] 250 katas covering core Zsh idioms and common anti-patterns.
- [x] 500 katas — halfway to 1.0.

### Version 1.0.0 — the 1000-kata milestone (complete)

- [x] 1000 katas (ZC1001–ZC2003) covering:
    - Syntax errors.
    - Portability issues.
    - Performance bottlenecks.
    - Security vulnerabilities.
    - Best practices.

### Version 1.x — beyond the milestone

- [ ] **Language Server Protocol (LSP).**
  An official LSP implementation for VS Code, Neovim, and other editors with inline diagnostics and quick-fix actions.

- [x] **Auto-fixer core.**
  `-fix`, `-diff`, and `-dry-run` flags apply deterministic per-kata rewrites since v1.0.14.
  The set of fix-enabled katas grows each release.

- [ ] **Auto-fixer coverage.**
  Expansion continues per release.
  The structural ceiling is the subset of detections that admit a context-free, idempotent, byte-exact rewrite.
  Many advisory or context-dependent detections remain detection-only by design.

- [ ] **Plugin system.**
  Custom checks authored in Lua or Wasm.

- [ ] **Distribution channels.**
  Broaden install paths beyond `install.sh`, `go install`, and the signed Releases archive:

  - [ ] **Homebrew (homebrew-core).**
    Submission planned once the [Homebrew notability threshold](https://docs.brew.sh/Acceptable-Formulae) is met (≥75 stars / ≥30 forks / ≥30 watchers for a third-party submission, or ≥225 / ≥90 / ≥90 for a self-submission).
    The `homebrew-eligibility` workflow files a tracking issue automatically the week the threshold is crossed.
    Community-filed third-party submissions face the lower bar.
  - [x] **Linux packages.**
    `.deb`, `.rpm`, and `.apk` produced via goreleaser `nfpms:` and shipped on every Releases tag.
  - [ ] **AUR.**
    `zshellcheck-bin` published via goreleaser `aurs:` for Arch and Manjaro.
  - [x] **PowerShell installer.**
    `install.ps1` mirrors `install.sh` for Windows.
  - [x] **Docker image.**
    `FROM scratch` static binary published to `ghcr.io/afadesigns/zshellcheck` on tag, signed by cosign.

## Long-term vision

- **Type checking.**
  Experimental static type inference for Zsh scripts.
- **Formatter.**
  A strictly opinionated formatter for Zsh, in the spirit of `gofmt` and `prettier`.

## Progress

The 1.0 milestone shipped 1000 katas.
v1.x hotfixes — parser, dedup, severity rebalance, XDG config, inline disable directives, Marketplace action rename — shipped between v1.0.0 and the latest tag.
See [CHANGELOG.md](CHANGELOG.md) for the full per-release history and [KATAS.md](KATAS.md) for the current kata list.
