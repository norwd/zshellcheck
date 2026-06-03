#!/usr/bin/env bash
# violation-corpus-sweep.sh — Run zshellcheck against the pinned upstream
# corpora and compare per-file, per-kata violation counts to a baseline.
#
# This is the false-positive ratchet. The parser-corpus sweep proves the
# corpora never crash and never produce parser errors; this sweep snapshots
# the kata findings on that same known-good code. Any drift — a finding that
# appears, disappears, or changes count — fails the gate and must be reviewed.
# A new finding on code that a human already reviewed as clean is a candidate
# false positive; a vanished finding is a candidate regression in coverage.
# The same snapshot-ratchet pattern backs shellcheck, Ruff, clippy, and semgrep.
#
# Manifest: .github/parser-corpus-manifest.txt (shared with the parser sweep)
# Baseline: .github/violation-baseline.txt
#
# Usage:
#   scripts/violation-corpus-sweep.sh                   # gate (compare to baseline)
#   scripts/violation-corpus-sweep.sh --update-baseline # rewrite baseline from current run
#   ZSHELLCHECK_BIN=/path/to/binary scripts/violation-corpus-sweep.sh
#
# Exit status:
#   0  every file matches its baseline
#   1  one or more findings drifted (review, then --update-baseline if intended)
#   2  manifest / clone / IO error, or a linter panic (never tolerated)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="${REPO_ROOT}/.github/parser-corpus-manifest.txt"
BASELINE="${REPO_ROOT}/.github/violation-baseline.txt"
WORK_DIR="${PARSER_SWEEP_WORK:-${REPO_ROOT}/testdata/external-corpora}"
BIN="${ZSHELLCHECK_BIN:-${REPO_ROOT}/zshellcheck}"

UPDATE_BASELINE=0
if [[ "${1:-}" == "--update-baseline" ]]; then
    UPDATE_BASELINE=1
fi

if [[ ! -f "${MANIFEST}" ]]; then
    echo "::error::manifest not found: ${MANIFEST}" >&2
    exit 2
fi

if [[ ! -x "${BIN}" ]]; then
    echo ">>> building zshellcheck into ${BIN}"
    (cd "${REPO_ROOT}" && go build -o "${BIN}" ./cmd/zshellcheck)
fi

mkdir -p "${WORK_DIR}"

# Fetch a corpus at a pinned SHA. Idempotent. Mirrors the parser sweep.
fetch_corpus() {
    local name="$1" sha="$2" url="$3"
    local dir="${WORK_DIR}/${name}"
    if [[ -d "${dir}/.git" ]]; then
        local cur
        cur=$(cd "${dir}" && git rev-parse HEAD 2>/dev/null || echo "")
        if [[ "${cur}" == "${sha}"* ]]; then
            return 0
        fi
        (cd "${dir}" && git fetch --quiet --depth=1 origin "${sha}" 2>/dev/null \
            && git reset --quiet --hard "${sha}") || {
            rm -rf "${dir}"
        }
    fi
    if [[ ! -d "${dir}/.git" ]]; then
        git clone --quiet "${url}" "${dir}"
        (cd "${dir}" && git reset --quiet --hard "${sha}") || {
            echo "::error::failed to checkout ${sha} in ${name}" >&2
            return 2
        }
    fi
}

# Build the file list for one corpus given its glob list.
list_files() {
    local dir="$1"
    shift
    local globs=("$@")
    local find_args=()
    local first=1
    for g in "${globs[@]}"; do
        if (( first )); then
            find_args+=(-name "${g}")
            first=0
        else
            find_args+=(-o -name "${g}")
        fi
    done
    find "${dir}" -type f \( "${find_args[@]}" \) -print 2>/dev/null | sort
}

declare -A CURRENT
declare -A BASELINE_MAP

if [[ -f "${BASELINE}" ]]; then
    while IFS=$'\t' read -r path kata count; do
        [[ -z "${path}" || "${path}" == \#* ]] && continue
        BASELINE_MAP["${path}	${kata}"]="${count}"
    done < "${BASELINE}"
fi

total_files=0
total_findings=0
PANICS=()

while IFS=$'\t' read -r name sha url glob_list; do
    [[ -z "${name}" || "${name}" == \#* ]] && continue
    fetch_corpus "${name}" "${sha}" "${url}" || exit 2

    # Split the glob list on spaces without pathname expansion (see the parser
    # sweep for the rationale: a bare expansion would let _* match repo-root
    # files before find ever sees it).
    read -r -a globs <<< "${glob_list}"
    while IFS= read -r f; do
        [[ -z "${f}" ]] && continue
        rel="${name}/${f#${WORK_DIR}/${name}/}"
        # Run once in the main shell so a panic can be recorded. A panic exits
        # >= 2 with a Go stack trace; it is always fatal, even under
        # --update-baseline, so a crash can never be baked into the baseline.
        set +e
        out=$("${BIN}" -no-banner -format json "${f}" 2>&1)
        rc=$?
        set -e
        if (( rc >= 2 )) || grep -qiE 'panic:|^goroutine [0-9]+ ' <<< "${out}"; then
            PANICS+=("${rel} (exit ${rc})")
        fi
        total_files=$((total_files+1))
        # Aggregate per-kata counts for this file from the JSON KataID keys.
        while IFS= read -r kata; do
            [[ -z "${kata}" ]] && continue
            CURRENT["${rel}	${kata}"]=$(( ${CURRENT["${rel}	${kata}"]:-0} + 1 ))
            total_findings=$((total_findings+1))
        done < <(grep -oE '"KataID": "ZC[0-9]+"' <<< "${out}" | grep -oE 'ZC[0-9]+')
    done < <(list_files "${WORK_DIR}/${name}" "${globs[@]}")
done < "${MANIFEST}"

echo "violation-corpus sweep: files=${total_files} findings=${total_findings} keys=${#CURRENT[@]}"

# A panic is always fatal — even in --update-baseline mode.
if (( ${#PANICS[@]} > 0 )); then
    echo "::error::violation-corpus gate: ${#PANICS[@]} file(s) panicked the linter — integrations must be panic-free"
    for p in "${PANICS[@]}"; do
        echo "  ! ${p}"
    done
    exit 2
fi

if (( UPDATE_BASELINE )); then
    {
        echo "# Violation baseline. Regenerate with scripts/violation-corpus-sweep.sh --update-baseline"
        echo "# Snapshot of kata findings on the pinned, zero-parser-error corpora. Any drift is reviewed."
        echo "# Format: <relpath>\t<ZCID>\t<count>"
        while IFS= read -r k; do
            printf '%s\t%s\n' "${k}" "${CURRENT[$k]}"
        done < <(printf '%s\n' "${!CURRENT[@]}" | sort)
    } > "${BASELINE}"
    echo "wrote baseline: ${BASELINE} (${#CURRENT[@]} keys)"
    exit 0
fi

# Symmetric diff between CURRENT and BASELINE_MAP.
appeared=()   # new or higher count — candidate false positives
vanished=()   # gone or lower count — candidate coverage regressions
for k in "${!CURRENT[@]}"; do
    cur="${CURRENT[$k]}"
    base="${BASELINE_MAP[$k]:-0}"
    if (( cur > base )); then
        appeared+=("${k}: ${base} -> ${cur}")
    fi
done
for k in "${!BASELINE_MAP[@]}"; do
    base="${BASELINE_MAP[$k]}"
    cur="${CURRENT[$k]:-0}"
    if (( cur < base )); then
        vanished+=("${k}: ${base} -> ${cur}")
    fi
done

if (( ${#appeared[@]} > 0 || ${#vanished[@]} > 0 )); then
    echo "::error::violation-corpus gate: findings drifted from the baseline. Review each line, then run scripts/violation-corpus-sweep.sh --update-baseline if the change is intended."
    if (( ${#appeared[@]} > 0 )); then
        echo "  appeared (candidate false positives):"
        printf '%s\n' "${appeared[@]}" | sort | while IFS= read -r a; do
            echo "    + ${a}"
        done
    fi
    if (( ${#vanished[@]} > 0 )); then
        echo "  vanished (candidate coverage regressions):"
        printf '%s\n' "${vanished[@]}" | sort | while IFS= read -r v; do
            echo "    - ${v}"
        done
    fi
    exit 1
fi

exit 0
