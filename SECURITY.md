# Security policy

## Supported versions

| Version | Supported |
| --- | --- |
| `v1.0.x` (latest minor) | Yes |
| `< v1.0.0` | No |

Only the latest `v1.0.x` release receives security fixes.
Upgrade to the latest tag before reporting; fixes land as new patch releases, not as backports.

## Reporting a vulnerability

If you discover a vulnerability in ZShellCheck, disclose it responsibly using the process below.

### Process

1. **Do not open a public GitHub issue.**
   Public disclosure before a fix lets the issue be exploited.
2. **Use one of the two private channels.**
   - GitHub Private Vulnerability Reporting (preferred): submit at [Security → Advisories → Report a vulnerability](https://github.com/afadesigns/zshellcheck/security/advisories/new).
     The form is encrypted in transit and visible only to maintainers.
   - Email the maintainer at `github@afadesign.co`.
     GitHub private contact is also available via [@afadesigns](https://github.com/afadesigns).
3. **Include as much detail as possible.**
   - The type of vulnerability.
   - Full reproduction steps.
   - Any special configuration required.
   - Potential impact.

### Response process

1. **Acknowledge** the report within 7 days.
2. **Triage**: confirm the bug, assess CVSS severity (CVSS 4.0), establish a patch plan.
3. **Coordinate** with the reporter on a disclosure timeline.
   Default 90 days for medium-severity findings; shorter for critical.
4. **Patch**: land the fix on a private branch, request a CVE if applicable, prepare a release.
5. **Release**: cut a patch version, ship signed binaries, publish a security advisory on the [GitHub Security Advisories](https://github.com/afadesigns/zshellcheck/security/advisories) page.
6. **Credit** the reporter in the advisory and CHANGELOG, unless anonymity was requested.
7. **Post-mortem**: file a public issue describing root cause and prevention work, after the embargo lifts.

ZShellCheck is maintained by a solo developer; critical issues are triaged sooner, but a same-day response is not guaranteed.

## Vulnerability categories

ZShellCheck is a static-analysis tool.
Vulnerabilities fall into three categories:

1. **Code execution.**
   A malicious Zsh script causing ZShellCheck to execute arbitrary code on the host running the linter.
2. **Denial of service.**
   A malicious Zsh script causing ZShellCheck to hang or crash.
3. **False negatives.**
   Failure to report a critical security flaw in a Zsh script — for example a missed `eval` or injection.
   This is a bug class; high-impact misses are treated with high priority.

## Security requirements

What users can expect from ZShellCheck:

- The CLI must read input files and never execute their contents.
  The Zsh source under analysis is parsed, not invoked.
- The CLI must terminate on every input, valid or malformed.
  Parser depth and time bounds prevent runaway recursion.
- The CLI must never write outside paths the user invoked it on, except for the `-fix` mode which writes back to the same input files.
- Released binaries must be signed with cosign keyless OIDC and verifiable against the certificate identity `https://github.com/afadesigns/zshellcheck/.*`.
- Releases must ship with SHA-256 checksums (`checksums.txt`) signed alongside the artefacts.
- Build provenance must be recorded as SLSA Level 3 attestations, queryable from the GitHub attestations index.

What users should not expect:

- ZShellCheck is not a sandbox.
  Running `-fix` on a Zsh file that has already been compromised does not undo the compromise.
- ZShellCheck does not detect every Zsh anti-pattern.
  Coverage is bounded by the kata catalog; missing detections are bug class 3 above, not vulnerabilities in the linter itself.

For threat model, trust boundaries, and the assurance argument see [docs/THREAT_MODEL.md](docs/THREAT_MODEL.md).
