package checks

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Result struct {
	Category    string
	Name        string
	Requirement string
	Score       float64
	Earned      float64
	Passed      bool
	Details     string
}

type RuntimeOptions struct {
	InstallDeps      bool
	RunCommands      bool
	RunE2E           bool
	RunOffline       bool
	RunAccessibility bool
	RunResponsive    bool
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
	score := scoreForCheck(name)
	return Result{Category: categoryForCheck(name), Name: name, Requirement: requirementForCheck(name), Score: score, Earned: score, Passed: true, Details: details}
}

func skip(name, details string) Result {
	return Result{Category: categoryForCheck(name), Name: name, Requirement: requirementForCheck(name), Score: scoreForCheck(name), Earned: 0, Passed: true, Details: details}
}

func fail(name, details string) Result {
	requirement := requirementForCheck(name)
	if requirement != "" && !strings.HasPrefix(details, "section ") {
		details = "section " + requirement + ": " + details
	}
	return Result{Category: categoryForCheck(name), Name: name, Requirement: requirement, Score: scoreForCheck(name), Earned: 0, Passed: false, Details: details}
}

func categoryForCheck(name string) string {
	categories := map[string]string{
		"spec checklist":            "meta",
		"required files":            "structure",
		"required stack":            "static contracts",
		"naming conventions":        "structure",
		"package scripts":           "static contracts",
		"no remote backend markers": "static contracts",
		"storage keys":              "static contracts",
		"local persistence usage":   "static contracts",
		"route contract files":      "static contracts",
		"splash timing":             "static contracts",
		"type contracts":            "static contracts",
		"utility contracts":         "static contracts",
		"ui contract markers":       "static contracts",
		"auth behavior markers":     "source behavior",
		"habit behavior markers":    "source behavior",
		"source behavior markers":   "source behavior",
		"accessibility markers":     "quality gates",
		"pwa contract":              "pwa",
		"test suite contract":       "quality gates",
		"coverage config":           "quality gates",
		"mentor ai check":           "quality gates",
		"readme requirements":       "documentation",
		"executable verification":   "runtime behavior",
	}
	if category, ok := categories[name]; ok {
		return category
	}
	return "uncategorized"
}

func requirementForCheck(name string) string {
	requirements := map[string]string{
		"required stack":            "3",
		"no remote backend markers": "3",
		"route contract files":      "4",
		"splash timing":             "4",
		"storage keys":              "5",
		"local persistence usage":   "5",
		"required files":            "6",
		"naming conventions":        "7",
		"type contracts":            "8",
		"utility contracts":         "9",
		"ui contract markers":       "10",
		"auth behavior markers":     "11",
		"habit behavior markers":    "12",
		"source behavior markers":   "11-12",
		"pwa contract":              "13",
		"accessibility markers":     "15",
		"test suite contract":       "16",
		"coverage config":           "17",
		"package scripts":           "18",
		"readme requirements":       "19",
		"executable verification":   "4,9,11-18",
	}
	return requirements[name]
}

func scoreForCheck(name string) float64 {
	scores := map[string]float64{
		"required stack":            0.75,
		"required files":            1.00,
		"package scripts":           0.50,
		"no remote backend markers": 0.25,
		"storage keys":              0.50,
		"local persistence usage":   0.50,
		"route contract files":      0.75,
		"splash timing":             0.25,
		"type contracts":            0.50,
		"utility contracts":         0.75,
		"naming conventions":        0.25,
		"ui contract markers":       0.75,
		"auth behavior markers":     0.50,
		"habit behavior markers":    0.50,
		"source behavior markers":   1.00,
		"accessibility markers":     0.75,
		"pwa contract":              0.75,
		"test suite contract":       1.00,
		"coverage config":           0.50,
		"readme requirements":       0.25,
		"executable verification":   1.00,
	}
	return scores[name]
}

func TotalScore(results []Result) (earned, possible float64) {
	for _, result := range results {
		possible += result.Score
		earned += result.Earned
	}
	return earned, possible
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
