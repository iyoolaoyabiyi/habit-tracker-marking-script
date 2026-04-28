# Habit Tracker Examiner

`habit-examiner` is a Go CLI for reviewing a submitted Habit Tracker PWA against the technical requirements document.

This is a gate check, if any required part of each check is missing, the contract is broken and the submission fails.

It can inspect either:

- a local repository path
- a GitHub repository URL

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

### Keep the temporary clone

```bash
./bin/habit-examiner -repo https://github.com/owner/repo -keep-clone
```

## Flags

- `-repo`
  Local repository path or GitHub repository URL.
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
