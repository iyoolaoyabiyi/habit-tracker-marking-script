package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"habit-tracker-examiner/examiner/checks"
)

func main() {
	repoPath := flag.String("repo", ".", "path to a local repository or a GitHub repository URL")
	keepClone := flag.Bool("keep-clone", false, "keep the temporary cloned repository when -repo is a URL")
	installDeps := flag.Bool("install-deps", true, "install JavaScript dependencies before running executable checks")
	runCommands := flag.Bool("run-commands", true, "run build and test scripts")
	runE2E := flag.Bool("run-e2e", true, "run the Playwright end-to-end suite during executable verification")
	runOffline := flag.Bool("run-offline", true, "run examiner-owned offline/PWA browser checks when E2E is enabled")
	runAccessibility := flag.Bool("run-accessibility", true, "run examiner-owned accessibility browser checks when E2E is enabled")
	runResponsive := flag.Bool("run-responsive", true, "run examiner-owned responsive layout browser checks when E2E is enabled")
	flag.Parse()

	root, cleanup, err := resolveRepo(*repoPath, *keepClone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve repo: %v\n", err)
		os.Exit(2)
	}
	if cleanup != nil {
		defer cleanup()
	}

	fmt.Printf("Examining repository: %s\n\n", root)

	options := checks.RuntimeOptions{
		InstallDeps:      *installDeps,
		RunCommands:      *runCommands,
		RunE2E:           *runE2E,
		RunOffline:       *runOffline,
		RunAccessibility: *runAccessibility,
		RunResponsive:    *runResponsive,
	}

	results := checks.RunAll(root, options)

	var failed int
	for _, result := range results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
			failed++
		}

		fmt.Printf("[%s] %s / %s\n", status, result.Category, result.Name)
		if result.Details != "" {
			fmt.Printf("       %s\n", result.Details)
		}
	}

	if failed > 0 {
		fmt.Printf("\n%d check group(s) failed\n", failed)
		os.Exit(1)
	}

	fmt.Println("\nAll examiner checks passed")
}

func resolveRepo(repoArg string, keepClone bool) (string, func(), error) {
	if looksLikeGitHubRepo(repoArg) {
		return cloneRepo(repoArg, keepClone)
	}

	root, err := filepath.Abs(repoArg)
	if err != nil {
		return "", nil, err
	}

	info, err := os.Stat(root)
	if err != nil {
		return "", nil, err
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf("%s is not a directory", root)
	}

	return root, nil, nil
}

func looksLikeGitHubRepo(value string) bool {
	return strings.HasPrefix(value, "https://github.com/") ||
		strings.HasPrefix(value, "http://github.com/") ||
		strings.HasPrefix(value, "git@github.com:")
}

func cloneRepo(repoURL string, keepClone bool) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "habit-tracker-examiner-*")
	if err != nil {
		return "", nil, err
	}

	cloneDir := filepath.Join(tempDir, "repo")
	fmt.Printf("Cloning repository %s into %s\n", repoURL, cloneDir)

	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", nil, fmt.Errorf("git clone failed: %w", err)
	}

	if keepClone {
		return cloneDir, func() {
			fmt.Printf("\nKept cloned repository at %s\n", cloneDir)
		}, nil
	}

	return cloneDir, func() {
		_ = os.RemoveAll(tempDir)
	}, nil
}

var pathPattern = regexp.MustCompile(`^[a-z0-9/_\-.]+$`)

func init() {
	if !pathPattern.MatchString("cmd/examiner/main.go") {
		panic(errors.New("unexpected path validation failure"))
	}
}
