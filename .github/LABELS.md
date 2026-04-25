# Labels

ZShellCheck issues and pull requests are categorised with the label set below.
Each family uses a shared color so the kanban reads at a glance, and the
prefix-based naming follows the convention popularised by Kubernetes and
Hashicorp ecosystem projects.

## Conventional commit types (no prefix)

These match the [Conventional Commits](https://www.conventionalcommits.org/)
spec and the commit grammar enforced by the `commitlint` pre-commit step.

| Label | Color | When |
| :--- | :--- | :--- |
| `feat` | `#a2eeef` | New feature or capability |
| `fix` | `#d73a4a` | Bug report or fix PR |
| `docs` | `#0075ca` | Documentation only |
| `ci` | `#c5def5` | CI / build / release workflow |
| `chore` | `#cfd3d7` | Maintenance, no behaviour change |
| `refactor` | `#fcfc00` | Refactor, no behaviour change |
| `test` | `#fbca04` | Tests added, fixed, or refactored |
| `perf` | `#bfdbfe` | Performance change |

## `area/` — codebase region

Each issue should carry **at most one** `area/*` label. The auto-labeler
applies these from changed paths in `.github/labeler.yml`.

| Label | Color | Covers |
| :--- | :--- | :--- |
| `area/parser` | `#5319e7` | `pkg/lexer`, `pkg/parser`, `pkg/ast` |
| `area/katas` | `#7044e8` | `pkg/katas` — detection logic |
| `area/cli` | `#8557e8` | `cmd/zshellcheck`, install scripts, man pages |
| `area/zsh` | `#a06ae9` | Zsh-language semantics: builtin, glob, option, arithmetic, parameter expansion |
| `area/source` | `#bf9bf3` | Generic source-tree change (fallback when no specific area fits) |

## `severity/` — kata severity

Tracks the severity level of a detection rule (matches the `Severity` field on
each kata in `pkg/katas/`). Apply at most one.

| Label | Color | Meaning |
| :--- | :--- | :--- |
| `severity/error` | `#b60205` | Breaks at runtime |
| `severity/warning` | `#d93f0b` | Dangerous behaviour |
| `severity/info` | `#fef2c0` | Edge case, fragile |
| `severity/style` | `#0e8a16` | Idiomatic convention |

## `status/` — triage state

A single status label may be applied at any point in an issue's life.

| Label | Color | Meaning |
| :--- | :--- | :--- |
| `status/needs-triage` | `#fbca04` | New issue/PR awaiting first triage (auto-applied by issue templates) |
| `status/needs-review` | `#bfbfbf` | Ready for maintainer review |

## Quality (kata feedback)

| Label | Color | Meaning |
| :--- | :--- | :--- |
| `false-positive` | `#e99695` | Kata fires when it shouldn't |
| `regression` | `#7a0000` | Used to work, now broken |
| `security` | `#3d0c11` | Security concern or advisory |

## Automation

| Label | Color | Source |
| :--- | :--- | :--- |
| `deps` | `#0366d6` | Auto-applied to dependabot PRs (`dependabot.yml`) |

## GitHub-discoverability defaults

Kept verbatim because GitHub surfaces them in cross-repo search and the
"good first issue" / "help wanted" tabs.

| Label | Color | Source |
| :--- | :--- | :--- |
| `good first issue` | `#7057ff` | GitHub default — newcomer-friendly |
| `help wanted` | `#008672` | GitHub default — contributor sought |
| `duplicate` | `#cfd3d7` | GitHub default |

## Application rules

- **One per family.** A single issue carries at most one `area/*`, one
  `severity/*`, one `status/*`. Conventional-type labels (`feat`/`fix`/...)
  also single-apply per issue.
- **Labels are renamed, never deleted.** Renames preserve the history of
  every attached issue and PR; deletes vacate that history.
- **Triage flow.** New issues land with `status/needs-triage` (auto-applied
  by issue templates). The maintainer removes that label after assigning the
  appropriate type / area / severity. Once the issue is being worked on, swap
  to `status/needs-review` when ready for review.
