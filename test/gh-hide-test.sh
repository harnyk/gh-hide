#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

mkdir -p "$tmp/bin"
export GH_LOG="$tmp/gh.log"
cat > "$tmp/bin/gh" <<'GH'
#!/usr/bin/env bash
set -euo pipefail

printf '%q ' "$@" >> "$GH_LOG"
printf '\n' >> "$GH_LOG"

case "$1 $2" in
  'pr view')
    if [[ " $* " == *' number '* ]]; then
      printf '42\n'
    else
      printf 'PR_node_42\n'
    fi
    ;;
  'api graphql')
    if [[ " $* " == *CountViewedFiles* ]]; then
      printf '%s\0' 'UNVIEWED' 'VIEWED' 'VIEWED' 'UNVIEWED' 'UNVIEWED'
    elif [[ " $* " == *ListPullRequestFiles* ]]; then
      printf '%s\0' \
        'src/widget.spec.ts' 'VIEWED' \
        'src/xspecy.ts' 'UNVIEWED' \
        'test name.spec.ts' 'UNVIEWED' \
        'line
break.spec.ts' 'UNVIEWED' \
        'src/widget.ts' 'UNVIEWED'
    elif [[ " $* " == *' query '* ]]; then
      printf 'PR_node_42\n'
    fi
    ;;
esac
GH
chmod +x "$tmp/bin/gh"

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

assert_contains() {
  local haystack=$1 needle=$2
  [[ $haystack == *"$needle"* ]] || fail "expected output to contain: $needle"
}

assert_not_contains() {
  local haystack=$1 needle=$2
  [[ $haystack != *"$needle"* ]] || fail "did not expect output to contain: $needle"
}

run_hide() {
  : > "$GH_LOG"
  set +e
  output=$(PATH="$tmp/bin:$PATH" GH_LOG="$GH_LOG" "$root/gh-hide" "$@" 2>&1)
  status=$?
  set -e
}

run_hide '.spec.'
[[ $status -eq 0 ]] || fail "normal run failed: $output"
assert_contains "$output" 'Found PR #42.'
assert_contains "$output" 'Marking test name.spec.ts as viewed...'
assert_not_contains "$output" 'Marking src/widget.spec.ts as viewed...'
assert_contains "$output" 'Marked 2 file(s) as viewed.'
assert_contains "$output" 'Skipping src/widget.spec.ts (already viewed).'
assert_contains "$output" 'Skipped 1 already viewed file(s).'
log=$(<"$GH_LOG")
assert_contains "$log" 'ListPullRequestFiles'
assert_contains "$log" 'test\ name.spec.ts'
assert_contains "$log" 'PR_node_42'
assert_not_contains "$log" 'src/widget.spec.ts\ --'
assert_not_contains "$log" 'src/xspecy.ts'
assert_not_contains "$log" 'src/widget.ts'

run_hide 'no-match'
[[ $status -eq 0 ]] || fail "no-match run failed: $output"
assert_contains "$output" 'Marked 0 file(s) as viewed.'
log=$(<"$GH_LOG")
assert_not_contains "$log" 'markFileAsViewed'

for option in --dry-run -n; do
  run_hide "$option" '.spec.'
  [[ $status -eq 0 ]] || fail "$option run failed: $output"
  assert_contains "$output" 'Found PR #42.'
  assert_not_contains "$output" 'Would mark src/widget.spec.ts as viewed.'
  assert_contains "$output" 'test name.spec.ts'
  assert_contains "$output" 'Would mark 2 file(s) as viewed.'
  assert_contains "$output" 'Skipping src/widget.spec.ts (already viewed).'
  assert_contains "$output" 'Skipped 1 already viewed file(s).'
  log=$(<"$GH_LOG")
  assert_not_contains "$log" 'markFileAsViewed'
done

for option in --count -c; do
  run_hide "$option"
  [[ $status -eq 0 ]] || fail "$option run failed: $output"
  assert_contains "$output" 'visible: 3'
  assert_contains "$output" 'hidden : 2'
  assert_contains "$output" 'total  : 5'
  log=$(<"$GH_LOG")
  assert_contains "$log" 'CountViewedFiles'
done

run_hide ''
[[ $status -ne 0 ]] || fail 'empty filter unexpectedly succeeded'
assert_contains "$output" 'Usage:'
[[ ! -s $GH_LOG ]] || fail 'empty filter made a gh request'

run_hide
[[ $status -ne 0 ]] || fail 'missing filter unexpectedly succeeded'
assert_contains "$output" 'Usage:'
[[ ! -s $GH_LOG ]] || fail 'missing filter made a gh request'

printf 'PASS: gh-hide\n'
