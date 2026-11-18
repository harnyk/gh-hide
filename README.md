# gh-hide

A [GitHub CLI](https://cli.github.com/) extension (`gh hide`) that marks pull request files as viewed so they collapse in the review UI and stay out of the way. The filter is a plain substring of the file path—not a regex.

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`), authenticated (`gh auth login`)
- Bash
- Current branch is checked out locally and has an open PR (`gh pr view` resolves it)

## Installation

### From GitHub

```bash
gh extension install harnyk/gh-hide
```

### From a local clone

```bash
cd /path/to/gh-hide
gh extension install .
```

Verify:

```bash
gh hide --help   # not always available; more reliable:
gh extension list | grep gh-hide
```

### Upgrade / remove

```bash
gh extension upgrade gh-hide
gh extension remove gh-hide
```

## Usage

Run commands from the **repository root on a branch with an open PR**.

### Hide files by path substring

A file matches when its PR path **contains** `FILTER` literally (substring match).

```bash
# anything with .spec. in the path
gh hide '.spec.'

# lockfiles
gh hide 'pnpm-lock.yaml'

# everything under a directory
gh hide 'generated/'
```

### Dry run (no API changes)

```bash
gh hide --dry-run '.spec.'
gh hide -n '.spec.'
```

### Hide every file in the PR

```bash
gh hide --all
gh hide -a
```

Do not pass a filter with `--all`.

### Unhide (mark as not viewed again)

```bash
gh hide --unhide '.spec.'
gh hide -u '.spec.'

# everything currently marked viewed
gh hide --all --unhide
```

### Count viewed vs not viewed

```bash
gh hide --count
gh hide -c
```

Example output:

```text
visible: 12
hidden : 45
total  : 57
```

## Behavior

- Files already marked viewed are skipped on a normal `hide` (message on stderr).
- On `--unhide`, files that are not viewed yet are skipped.
- No matches → exit code 0 and `Marked 0 file(s) as viewed.`
- `gh` errors or invalid arguments → non-zero exit code.

Usage (built into the extension):

```bash
gh hide
```

## Development and tests

Tests stub `gh`; no live API calls:

```bash
bash test/gh-hide-test.sh
```

Expected: `PASS: gh-hide`.