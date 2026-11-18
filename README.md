# gh-hide

A [GitHub CLI](https://cli.github.com/) extension (`gh hide`) that marks pull request files as viewed so they collapse in the review UI and stay out of the way. The filter is a plain substring of the file path—not a regex.

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`), authenticated (`gh auth login`)
- Current branch is checked out locally and has an open PR (`gh pr view` resolves it)

## Installation

### From GitHub

```bash
gh extension install harnyk/gh-hide
```

`gh` downloads the precompiled binary for your platform from the latest GitHub release.

### From a local clone

Requires Go (see `go.mod` for the version).

```bash
cd /path/to/gh-hide
go build -o gh-hide .
gh extension install .
```

Rebuild with `go build -o gh-hide .` after pulling changes.

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

### Concurrency

Files are marked (or unmarked) in parallel, 8 at a time by default:

```bash
gh hide --concurrency 4 '.spec.'
gh hide -j 1 '.spec.'   # sequential
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
- Progress messages may appear in any order because files are processed concurrently.
- If a request fails, remaining requests are cancelled and the command exits with code 1.
- Invalid arguments → exit code 2; `gh`/API errors → exit code 1.

Built-in help:

```bash
gh hide --help
```

## Development and tests

Tests use a fake GitHub client; no live API calls:

```bash
go test -race ./...
```

## Releasing

Releases are built by [GoReleaser](https://goreleaser.com/) in GitHub Actions when a `v*` tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Check the configuration and build locally without publishing:

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```