package checks

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result struct {
	Name    string
	Passed  bool
	Details string
}

type RuntimeOptions struct {
	InstallDeps bool
	RunCommands bool
	RunE2E      bool
}

type manifest struct {
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	StartURL        string `json:"start_url"`
	Display         string `json:"display"`
	BackgroundColor string `json:"background_color"`
	ThemeColor      string `json:"theme_color"`
	Icons           []struct {
		Src   string `json:"src"`
		Sizes string `json:"sizes"`
		Type  string `json:"type"`
	} `json:"icons"`
}

func RunAll(root string, options RuntimeOptions) []Result {
	return []Result{
		checkSpecChecklist(),
		checkRequiredFiles(root),
		checkRequiredStack(root),
		checkNamingConventions(root),
		checkPackageScripts(root),
		checkNoRemoteBackendMarkers(root),
		checkStorageKeys(root),
		checkLocalPersistenceUsage(root),
		checkRouteFiles(root),
		checkSplashTiming(root),
		checkTypeContracts(root),
		checkUtilityContracts(root),
		checkUIContracts(root),
		checkAuthBehaviorMarkers(root),
		checkHabitBehaviorMarkers(root),
		checkSourceBehaviorMarkers(root),
		checkAccessibilityMarkers(root),
		checkPWAContracts(root),
		checkTests(root),
		checkCoverageConfig(root),
		checkMentorAICheck(root),
		checkReadme(root),
		checkExecutableVerification(root, options),
	}
}

func pass(name, details string) Result {
	return Result{Name: name, Passed: true, Details: details}
}

func fail(name, details string) Result {
	return Result{Name: name, Passed: false, Details: details}
}

func runCommand(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func readExistingFiles(root string, files []string) (string, error) {
	var builder strings.Builder
	for _, rel := range files {
		content, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return "", err
		}
		builder.Write(content)
		builder.WriteByte('\n')
	}
	return builder.String(), nil
}

func tail(output string, lineCount int) string {
	lines := strings.Split(output, "\n")
	if len(lines) <= lineCount {
		return output
	}
	return strings.Join(lines[len(lines)-lineCount:], "\n")
}
