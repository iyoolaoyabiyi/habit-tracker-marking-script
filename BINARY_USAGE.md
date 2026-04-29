# Habit Examiner Binary Usage

This guide is for using a prebuilt `habit-examiner` binary. You do not need Go or the source code unless you want to rebuild the executable.

## What The Binary Does

`habit-examiner` checks Habit Tracker PWA submissions against the project requirements. It can examine:

- local repository directories
- GitHub repository URLs
- a batch of repositories/directories listed in text files

The binary prints a report to the terminal and can also write the same report to a log file.

## Requirements

Install these on the machine running the binary:

- `git`, required for GitHub repository URLs
- `npm`, required when executable checks are enabled
- Playwright browser dependencies, required when E2E/browser checks are enabled

Go is not required when using the binary.

## Install Requirements

### macOS

Install Homebrew first if it is not already installed:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Install Git and Node.js/npm:

```bash
brew install git node
```

Install Playwright browser dependencies:

```bash
npx playwright install --with-deps
```

If `--with-deps` is not supported on your macOS setup, run:

```bash
npx playwright install
```

### Ubuntu/Debian

Install Git and Node.js/npm:

```bash
sudo apt update
sudo apt install -y git nodejs npm
```

Install Playwright browser dependencies:

```bash
npx playwright install --with-deps
```

If the command asks for `sudo`, rerun it with:

```bash
sudo npx playwright install --with-deps
```

### Windows

Install Git:

```powershell
winget install --id Git.Git -e
```

Install Node.js/npm:

```powershell
winget install --id OpenJS.NodeJS.LTS -e
```

Open a new PowerShell window, then install Playwright browsers:

```powershell
npx playwright install
```

### Verify Installation

Run:

```bash
git --version
npm --version
npx playwright --version
```

If you only run static checks with `-run-commands=false`, `npm` and Playwright are not required.

## Quick Start

Check one local repository:

```bash
./habit-examiner -repo /path/to/submission
```

Check one GitHub repository:

```bash
./habit-examiner -repo https://github.com/owner/repo
```

Check many repositories from a file and write a log:

```bash
./habit-examiner -repo-list repos.txt -log logs/examiner.log
```

Run static checks only:

```bash
./habit-examiner -repo-list repos.txt -run-commands=false -log logs/static-review.log
```

## Batch Inputs

For small batches, repeat `-repo` or `-dir`:

```bash
./habit-examiner -repo /path/to/a -repo /path/to/b -dir /path/to/c
```

For larger batches, use list files:

```bash
./habit-examiner -repo-list repos.txt -dir-list dirs.txt -log logs/examiner.log
```

List files are newline-delimited:

```text
# Blank lines and comments are ignored
/path/to/submission-a
/path/to/submission-b
https://github.com/owner/repo
```

Use `-repo-list` for local repository paths and GitHub URLs. Use `-dir-list` for local repository directories only.

## Logs

Use `-log` to save the full terminal report:

```bash
./habit-examiner -repo-list repos.txt -log logs/log.txt
```

The binary creates the log file and parent directory if needed. The same output still appears in the terminal.

## Dependency Caching

Executable checks install dependencies in each target repository by default. To avoid downloading the same packages repeatedly across a batch, the binary uses a shared npm cache.

Default cache location:

```text
/tmp/habit-examiner-npm-cache
```

Custom cache location:

```bash
./habit-examiner -repo-list repos.txt -npm-cache /tmp/habit-examiner-npm-cache
```

The cache reuses downloaded package tarballs. It does not share one `node_modules` directory across repositories. Each repository still gets its own `node_modules` tree based on its own lockfile, which keeps results correct when dependencies differ.

## Flags

- `-repo`
  Local repository path or GitHub repository URL. May be repeated.
- `-dir`
  Local repository directory path. May be repeated.
- `-repo-list`
  Newline-delimited file containing local repository paths or GitHub repository URLs.
- `-dir-list`
  Newline-delimited file containing local repository directory paths.
- `-log`
  Write output to a log file while still printing to the terminal.
- `-npm-cache`
  Shared npm cache directory for executable checks. Default: OS temp directory plus `habit-examiner-npm-cache`.
- `-keep-clone`
  Keep the temporary cloned repository when `-repo` is a GitHub URL.
- `-install-deps`
  Install JavaScript dependencies before executable checks. Default: `true`.
- `-run-commands`
  Run build and test scripts. Default: `true`.
- `-run-e2e`
  Run the submitted Playwright E2E suite during executable verification. Default: `true`.
- `-run-offline`
  Run examiner-owned offline/PWA browser checks when E2E is enabled. Default: `true`.
- `-run-accessibility`
  Run examiner-owned accessibility browser checks when E2E is enabled. Default: `true`.
- `-run-responsive`
  Run examiner-owned responsive layout checks when E2E is enabled. Default: `true`.

## Recommended Commands

Full review for many repositories:

```bash
./habit-examiner -repo-list repos.txt -log logs/full-review.log
```

Static-only review:

```bash
./habit-examiner -repo-list repos.txt -run-commands=false -log logs/static-review.log
```

Build, unit, and integration checks without browser/E2E checks:

```bash
./habit-examiner -repo-list repos.txt -run-e2e=false -log logs/no-e2e-review.log
```

Use existing `node_modules` in each target:

```bash
./habit-examiner -repo-list repos.txt -install-deps=false -log logs/no-install-review.log
```

Keep cloned GitHub repositories for inspection:

```bash
./habit-examiner -repo-list repos.txt -keep-clone -log logs/review.log
```

## What To Expect

For each target, the report starts with the requested target and resolved repository path:

```text
Examining repo: /path/to/submission
Resolved repository: /path/to/submission
```

Each check group prints a `PASS` or `FAIL` line:

```text
[PASS] structure / required files (1.00/1.00)
[FAIL] quality gates / mentor ai check (-1.00/0.00)
```

Failures include details indented below the failed check:

```text
[FAIL] static contracts / package scripts (0.00/0.50)
       section 18: Issues found:
       - package.json is missing test:e2e
```

Each target ends with a score:

```text
Score: 11.00/13.00
1 check group(s) failed
```

When multiple targets are provided, the binary continues through all targets even if one fails. The final process exit code reflects the overall result.

## Exit Codes

- `0`: all checked targets passed
- `1`: one or more targets completed but failed examiner checks
- `2`: one or more targets could not be resolved or set up

## Notes

- GitHub URL targets are cloned into a temporary directory and removed after the run unless `-keep-clone` is set.
- `npm ci` is used when `package-lock.json` exists. Otherwise the binary falls back to `npm install`.
- Static checks are fast and do not execute submitted app commands.
- Runtime and browser checks can take significantly longer because they install dependencies, build the app, run submitted tests, and run examiner-owned tests.
