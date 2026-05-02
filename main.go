package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"habit-tracker-examiner/checks"
)

type stringList []string

func (values *stringList) String() string {
	return strings.Join(*values, ", ")
}

func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func main() {
	var repoPaths stringList
	var dirPaths stringList
	flag.Var(&repoPaths, "repo", "path to a local repository or a GitHub repository URL; may be repeated")
	flag.Var(&dirPaths, "dir", "path to a local repository directory; may be repeated")
	repoListPath := flag.String("repo-list", "", "newline-delimited file of local repository paths or GitHub repository URLs")
	dirListPath := flag.String("dir-list", "", "newline-delimited file of local repository directories")
	logPath := flag.String("log", "", "write examiner output to this log file")
	npmCachePath := flag.String("npm-cache", filepath.Join(os.TempDir(), "habit-examiner-npm-cache"), "shared npm cache directory for executable checks")
	keepClone := flag.Bool("keep-clone", false, "keep the temporary cloned repository when -repo is a URL")
	installDeps := flag.Bool("install-deps", true, "install JavaScript dependencies before running executable checks")
	runCommands := flag.Bool("run-commands", true, "run build and test scripts")
	runE2E := flag.Bool("run-e2e", true, "run the Playwright end-to-end suite during executable verification")
	runOffline := flag.Bool("run-offline", true, "run examiner-owned offline/PWA browser checks when E2E is enabled")
	runAccessibility := flag.Bool("run-accessibility", true, "run examiner-owned accessibility browser checks when E2E is enabled")
	runResponsive := flag.Bool("run-responsive", true, "run examiner-owned responsive layout browser checks when E2E is enabled")
	flag.Parse()

	if *repoListPath != "" {
		values, err := readTargetList(*repoListPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read repo list: %v\n", err)
			os.Exit(2)
		}
		repoPaths = append(repoPaths, values...)
	}

	if *dirListPath != "" {
		values, err := readTargetList(*dirListPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read dir list: %v\n", err)
			os.Exit(2)
		}
		dirPaths = append(dirPaths, values...)
	}

	if len(repoPaths) == 0 && len(dirPaths) == 0 {
		repoPaths = append(repoPaths, ".")
	}

	writer := io.Writer(os.Stdout)
	var logFile *os.File
	if *logPath != "" {
		var err error
		if err := os.MkdirAll(filepath.Dir(*logPath), 0o755); err != nil && filepath.Dir(*logPath) != "." {
			fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
			os.Exit(2)
		}
		logFile, err = os.Create(*logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create log file: %v\n", err)
			os.Exit(2)
		}
		defer logFile.Close()
		writer = io.MultiWriter(os.Stdout, logFile)
	}

	options := checks.RuntimeOptions{
		InstallDeps:      *installDeps,
		RunCommands:      *runCommands,
		RunE2E:           *runE2E,
		RunOffline:       *runOffline,
		RunAccessibility: *runAccessibility,
		RunResponsive:    *runResponsive,
		NPMCache:         *npmCachePath,
	}

	var checkFailures int
	var setupFailures int

	for index, repoPath := range repoPaths {
		if index > 0 {
			fmt.Fprintln(writer)
		}
		switch runTarget(writer, "repo", repoPath, *keepClone, options) {
		case 1:
			checkFailures++
		case 2:
			setupFailures++
		}
	}

	for index, dirPath := range dirPaths {
		if len(repoPaths) > 0 || index > 0 {
			fmt.Fprintln(writer)
		}
		switch runTarget(writer, "dir", dirPath, false, options) {
		case 1:
			checkFailures++
		case 2:
			setupFailures++
		}
	}

	if setupFailures > 0 {
		os.Exit(2)
	}
	if checkFailures > 0 {
		os.Exit(1)
	}

	fmt.Fprintln(writer, "\nAll examiner targets passed")
}

func readTargetList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value == "" || strings.HasPrefix(value, "#") {
			continue
		}
		values = append(values, value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func runTarget(writer io.Writer, targetType, targetPath string, keepClone bool, options checks.RuntimeOptions) int {
	root, cleanup, err := resolveRepo(writer, targetPath, keepClone)
	if err != nil {
		fmt.Fprintf(writer, "[SETUP FAIL] %s:%s\n", targetType, targetPath)
		fmt.Fprintf(writer, "failed to resolve repo: %v\n", err)
		return 2
	}
	if cleanup != nil {
		defer cleanup()
	}

	fmt.Fprintf(writer, "Examining %s: %s\n", targetType, targetPath)
	fmt.Fprintf(writer, "Resolved repository: %s\n\n", root)

	results := checks.RunAll(root, options)

	var failed int
	for _, result := range results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
			failed++
		}

		fmt.Fprintf(writer, "[%s] %s / %s (%.2f/%.2f)\n", status, result.Category, result.Name, result.Earned, result.Score)
		if result.Details != "" {
			printIndented(writer, result.Details)
		}
	}
	earned, possible := checks.TotalScore(results)

	fmt.Fprintf(writer, "\nScore: %.2f/%.2f\n", earned, possible)
	if failed > 0 {
		fmt.Fprintf(writer, "%d check group(s) failed\n", failed)
		return 1
	}

	fmt.Fprintln(writer, "All examiner checks passed")
	return 0
}

func printIndented(writer io.Writer, details string) {
	for _, line := range strings.Split(details, "\n") {
		fmt.Fprintf(writer, "       %s\n", line)
	}
}

func resolveRepo(writer io.Writer, repoArg string, keepClone bool) (string, func(), error) {
	if target, ok, err := parseGitHubTreeURL(repoArg); ok || err != nil {
		if err != nil {
			return "", nil, err
		}
		return cloneRepoTarget(writer, target, keepClone)
	}

	if looksLikeGitHubRepo(repoArg) {
		return cloneRepoTarget(writer, githubTarget{CloneURL: repoArg}, keepClone)
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

type githubTarget struct {
	CloneURL string
	Branch   string
	Subdir   string
}

func parseGitHubTreeURL(value string) (githubTarget, bool, error) {
	parsed, err := url.Parse(value)
	if err != nil {
		return githubTarget{}, false, nil
	}
	if parsed.Host != "github.com" && parsed.Host != "www.github.com" {
		return githubTarget{}, false, nil
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "tree" {
		return githubTarget{}, false, nil
	}

	owner := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")
	if owner == "" || repo == "" {
		return githubTarget{}, true, fmt.Errorf("invalid GitHub tree URL: %s", value)
	}

	cloneURL := fmt.Sprintf("%s://github.com/%s/%s.git", parsed.Scheme, owner, repo)
	remainder := strings.Join(parts[3:], "/")
	branch, subdir, err := resolveGitHubTreeBranch(cloneURL, remainder)
	if err != nil {
		return githubTarget{}, true, err
	}

	return githubTarget{CloneURL: cloneURL, Branch: branch, Subdir: subdir}, true, nil
}

func resolveGitHubTreeBranch(cloneURL, remainder string) (string, string, error) {
	output, err := runLocalCommand("git", "ls-remote", "--heads", cloneURL)
	if err != nil {
		return "", "", fmt.Errorf("could not inspect remote branches for %s: %w", cloneURL, err)
	}

	branch, subdir := matchGitHubTreeBranch(output, remainder)
	if branch == "" {
		return "", "", fmt.Errorf("could not match GitHub tree branch %q in %s", remainder, cloneURL)
	}

	return branch, subdir, nil
}

func matchGitHubTreeBranch(remoteRefs, remainder string) (string, string) {
	var branch string
	for _, line := range strings.Split(remoteRefs, "\n") {
		_, ref, ok := strings.Cut(line, "\t")
		if !ok {
			continue
		}
		name := strings.TrimPrefix(ref, "refs/heads/")
		if remainder == name || strings.HasPrefix(remainder, name+"/") {
			if len(name) > len(branch) {
				branch = name
			}
		}
	}

	if branch == "" {
		return "", ""
	}

	if remainder == branch {
		return branch, ""
	}

	return branch, strings.TrimPrefix(remainder, branch+"/")
}

func runLocalCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func cloneRepoTarget(writer io.Writer, target githubTarget, keepClone bool) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "habit-tracker-examiner-*")
	if err != nil {
		return "", nil, err
	}

	cloneDir := filepath.Join(tempDir, "repo")
	fmt.Fprintf(writer, "Cloning repository %s into %s\n", target.CloneURL, cloneDir)

	args := []string{"clone", "--depth", "1"}
	if target.Branch != "" {
		args = append(args, "--branch", target.Branch)
	}
	args = append(args, target.CloneURL, cloneDir)

	cmd := exec.Command("git", args...)
	cmd.Stdout = writer
	cmd.Stderr = writer

	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", nil, fmt.Errorf("git clone failed: %w", err)
	}

	root := cloneDir
	if target.Subdir != "" {
		root = filepath.Join(cloneDir, filepath.FromSlash(target.Subdir))
		info, err := os.Stat(root)
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return "", nil, fmt.Errorf("GitHub tree subdirectory %s was not found after clone: %w", target.Subdir, err)
		}
		if !info.IsDir() {
			_ = os.RemoveAll(tempDir)
			return "", nil, fmt.Errorf("GitHub tree path %s is not a directory", target.Subdir)
		}
	}

	if keepClone {
		return root, func() {
			fmt.Fprintf(writer, "\nKept cloned repository at %s\n", cloneDir)
		}, nil
	}

	return root, func() {
		_ = os.RemoveAll(tempDir)
	}, nil
}

var pathPattern = regexp.MustCompile(`^[a-z0-9/_\-.]+$`)

func init() {
	if !pathPattern.MatchString("cmd/main.go") {
		panic(errors.New("unexpected path validation failure"))
	}
}
