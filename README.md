# Habit Tracker Examiner

`habit-examiner` is a Go CLI for reviewing a submitted Habit Tracker PWA against the technical requirements document.

This is a gate check, if any required part of each check is missing, the contract is broken and the submission fails.

It can inspect either:

- one or more local repository paths
- one or more GitHub repository URLs

## Setup And Use

From the `cmd` directory, verify Go is available:

```bash
go version
```

No separate Go dependency installation step is required. The examiner is a Go module and can be run directly from source:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo /path/to/submission
# or 
go run . -repo /path/to/submission
```

For a static-only review that does not install or execute the submitted app:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo /path/to/submission -run-commands=false
```

To examine a GitHub repository directly:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo https://github.com/owner/repo
```

To examine a subdirectory from a GitHub branch:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo https://github.com/owner/repo/tree/main/frontend/stage-3
```

To examine multiple local directories or repositories and save the report:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo /path/to/submission-a -dir /path/to/submission-b -log logs/examiner.log
```

For larger batches, use newline-delimited list files:

```bash
GOCACHE=/tmp/go-build-cache go run . -repo-list repos.txt -dir-list dirs.txt -log logs/examiner.log
```

Executable checks use a shared npm cache by default at `/tmp/habit-examiner-npm-cache` on Unix-like systems. This lets later repositories reuse package downloads while still installing a correct per-repository `node_modules` tree from each repository's lockfile.

## Generate A Binary

Build a reusable local binary:

```bash
mkdir -p bin
GOCACHE=/tmp/go-build-cache go build -o bin/habit-examiner .
```

Run the binary:

```bash
./bin/habit-examiner -repo /path/to/submission
```

You can also use the Makefile:

```bash
make build
./bin/habit-examiner -repo /path/to/submission
```

## What It Checks

- required files and route files
- required stack declaration and usage
- naming conventions
- required package scripts
- local `localStorage` persistence markers
- route, splash, auth, habit, accessibility, and PWA markers
- required test files and exact test titles
- Vitest coverage configuration for `src/lib`
- runtime coverage summary parsing for the 80% `src/lib` line threshold
- README coverage
- mentor hidden-marker AI audit
- optional executable verification by running build, submitted tests, and examiner-owned runtime/browser tests

## Prerequisites

- Go 1.25+
- `git`
- `npm` for executable verification
- Playwright browser dependencies if you run E2E checks

## Usage

### Check the current repository

```bash
./bin/habit-examiner -repo .
```

### Check a GitHub repository

```bash
./bin/habit-examiner -repo https://github.com/owner/repo
```

GitHub `tree` URLs are supported when the submission lives in a repository subdirectory:

```bash
./bin/habit-examiner -repo https://github.com/owner/repo/tree/main/frontend/stage-3
```

### Check multiple targets and write a log

```bash
./bin/habit-examiner -repo /path/to/submission-a -repo https://github.com/owner/repo -dir /path/to/submission-b -log logs/examiner.log
```

For large batches:

```bash
./bin/habit-examiner -repo-list repos.txt -dir-list dirs.txt -log logs/examiner.log
```

List files are newline-delimited. Blank lines and lines starting with `#` are ignored.

### Keep the temporary clone

```bash
./bin/habit-examiner -repo https://github.com/owner/repo -keep-clone
```

## Flags

- `-repo`
  Local repository path, GitHub repository URL, or GitHub `tree` URL. May be repeated.
- `-dir`
  Local repository directory path. May be repeated.
- `-repo-list`
  Newline-delimited file of local repository paths or GitHub repository URLs.
- `-dir-list`
  Newline-delimited file of local repository directory paths.
- `-log`
  Write the examiner output to a log file while still printing it to the terminal.
- `-npm-cache`
  Shared npm cache directory used by executable checks. Default: OS temp directory plus `habit-examiner-npm-cache`.
- `-keep-clone`
  Keep the temporary cloned repository when `-repo` is a URL.
- `-install-deps`
  Install JavaScript dependencies before executable checks. Default: `true`.
- `-run-commands`
  Run build and test scripts. Default: `true`.
- `-run-e2e`
  Run the Playwright E2E suite during executable verification. Default: `true`.
- `-run-offline`
  Run examiner-owned offline/PWA browser checks when E2E is enabled. Default: `true`.
- `-run-accessibility`
  Run examiner-owned accessibility browser checks when E2E is enabled. Default: `true`.
- `-run-responsive`
  Run examiner-owned responsive layout browser checks when E2E is enabled. Default: `true`.

## Recommended Commands

### Static-only review

Useful when you only want contract validation without running app commands:

```bash
./bin/habit-examiner -repo https://github.com/owner/repo -run-commands=false
```

### Build, unit, and integration checks only

Useful in restricted environments where Playwright web servers cannot start:

```bash
./bin/habit-examiner -repo https://github.com/owner/repo -run-e2e=false
```

### Reuse existing dependencies

Useful when the submission already has `node_modules` installed:

```bash
./bin/habit-examiner -repo /path/to/submission -install-deps=false
```

### Reuse package downloads across a batch

Useful when checking many submissions that share similar dependencies:

```bash
./bin/habit-examiner -repo-list repos.txt -npm-cache /tmp/habit-examiner-npm-cache
```

The examiner intentionally does not share one `node_modules` directory across repositories. That can produce incorrect results when lockfiles, native packages, or postinstall scripts differ. The shared npm cache keeps downloads efficient while each target remains isolated.

## Extract Scores

Use `extract_scores.py` to convert a full examiner log into a compact `repo: score` file.

Basic usage:

```bash
python3 extract_scores.py logs/101-176-scores.txt
```

By default, the script writes beside the input log using the suffix `-extracted.txt`. For example, `logs/101-176-scores.txt` produces:

```text
logs/101-176-scores-extracted.txt
```

To preserve the original repo-list order and include repositories that failed before scoring, pass the matching repo list:

```bash
python3 extract_scores.py logs/101-176-scores.txt repo/101-176.txt
```

When a repo list is provided:

- output follows the repo-list order
- setup failures are scored as `0`
- repos in the list that never received a score are scored as `0`
- markdown links such as `[url](url)` are normalized to the URL
- duplicate repo entries are kept in order

Choose a specific output file with `-o`:

```bash
python3 extract_scores.py logs/101-176-scores.txt repo/101-176.txt -o logs/extracts.txt
```

Append a new batch to an existing output file with `--append` or `-a`:

```bash
python3 extract_scores.py logs/101-176-scores.txt repo/101-176.txt -o logs/extracts.txt --append
```

Append mode adds a blank separator line before the new batch when the output file already has content.

The output format is:

```text
https://github.com/owner/repo: 7.25
https://github.com/owner/broken-repo: 0
```

## Exit Codes

- `0`: all checks passed
- `1`: one or more checks failed
- `2`: repository resolution or setup failure

## Output Shape

The tool prints one line per check group:

```text
[PASS] structure / required files
[FAIL] quality gates / mentor ai check
```

Each failure includes the relevant requirements section where one is available and all detected issues for that check group, so an examiner can see what broke without reading the code first and an intern can fix related problems in one pass.

## Scoring

The examiner reports a score for each check and a final total out of 13 points. Passing a scored check earns its full weight; failing a scored check earns zero for that check. The mentor-marker audit remains outside the possible-point total, but a failure carries a negative earned score.

## Enforcement Coverage

The examiner groups checks into:

- `structure`
- `static contracts`
- `source behavior`
- `runtime behavior`
- `pwa`
- `quality gates`
- `documentation`
- `meta`

Static checks enforce required paths, package markers, exported contracts, test IDs, README sections, and source behavior evidence. Static and content checks are strict but do not intentionally fail fast: they collect every missing marker they can find within the group, then return one binary pass/fail result. Runtime checks still fail fast on blocking setup or execution failures such as install, build, or server startup errors. Browser checks run examiner-owned Playwright flows for auth, habit CRUD, user isolation, PWA/offline behavior, accessibility, and responsive overflow.

Failure details use this shape:

```text
[FAIL] quality gates / accessibility markers (0.00/0.75)
       section 15: Issues found:
       - src/components/auth/LoginForm.tsx is missing htmlFor associations
       - src/components/habits/HabitForm.tsx is missing button elements
       - src/app/globals.css is missing :focus-visible styling
```

## Fixture Coverage

The Go test suite includes fixture-style tests for examiner behavior, including shallow submitted-test detection and coverage-summary parsing. Run them with:

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```

## Remaining Manual Review

The examiner is intentionally stricter than simple string scanning, but it still cannot fully judge product quality. Human review should still inspect visual polish, copy quality, overall UX, code maintainability, and whether the implementation is genuinely understandable rather than merely shaped to satisfy automated checks.

## Cross-Platform Builds

Build release binaries from the `cmd` directory:

```bash
mkdir -p bin

GOOS=linux GOARCH=amd64   GOCACHE=/tmp/go-build-cache go build -o bin/habit-examiner-linux-amd64 .
GOOS=darwin GOARCH=arm64  GOCACHE=/tmp/go-build-cache go build -o bin/habit-examiner-darwin-arm64 .
GOOS=windows GOARCH=amd64 GOCACHE=/tmp/go-build-cache go build -o bin/habit-examiner-windows-amd64.exe .
```

## Notes

- The mentor AI check is designed to fail when the hidden mentor marker appears in the submission.
- Full executable verification depends on the target environment being able to run the app toolchain.
- For GitHub URLs, the tool clones to a temporary directory and removes it unless `-keep-clone` is set.
