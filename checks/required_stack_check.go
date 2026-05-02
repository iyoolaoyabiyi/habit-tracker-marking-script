package checks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkRequiredStack(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return failWithIssues("required stack", []string{"package.json could not be read: " + err.Error()})
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fail("required stack", "invalid package.json: "+err.Error())
	}

	hasPackage := func(name string) bool {
		if _, ok := pkg.Dependencies[name]; ok {
			return true
		}
		if _, ok := pkg.DevDependencies[name]; ok {
			return true
		}
		return false
	}

	requiredPackages := []string{
		"next",
		"react",
		"react-dom",
		"typescript",
		"tailwindcss",
		"@playwright/test",
		"vitest",
		"@testing-library/react",
	}
	var issues []string
	for _, name := range requiredPackages {
		if !hasPackage(name) {
			issues = append(issues, "package.json missing package "+name)
		}
	}

	stackUsageChecks := []struct {
		label        string
		file         string
		patternSets  [][]string
		missingLabel string
	}{
		{
			label:       "Next.js App Router",
			file:        "src/app/layout.tsx",
			patternSets: [][]string{{"export default function RootLayout"}},
		},
		{
			label:       "React",
			file:        "src/components/auth/LoginForm.tsx",
			patternSets: [][]string{{"useState"}},
		},
		{
			label:       "TypeScript",
			file:        "tsconfig.json",
			patternSets: [][]string{{`"compilerOptions"`}},
		},
		{
			label:        "Tailwind CSS",
			file:         "src/app/globals.css",
			patternSets:  [][]string{{`@import "tailwindcss"`}, {`@tailwind base`, `@tailwind components`, `@tailwind utilities`}},
			missingLabel: `@import "tailwindcss" or Tailwind v3 @tailwind directives`,
		},
		{
			label:       "Playwright",
			file:        "playwright.config.ts",
			patternSets: [][]string{{"@playwright/test", "defineConfig"}},
		},
		{
			label:       "Vitest",
			file:        "vitest.config.ts",
			patternSets: [][]string{{"vitest/config", "defineConfig"}},
		},
		{
			label:       "React Testing Library",
			file:        "tests/integration/auth-flow.test.tsx",
			patternSets: [][]string{{"@testing-library/react"}},
		},
	}

	for _, check := range stackUsageChecks {
		fileContent, err := os.ReadFile(filepath.Join(root, check.file))
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s evidence missing: %s could not be read: %v", check.label, check.file, err))
			continue
		}
		text := string(fileContent)
		matched := false
		for _, patternSet := range check.patternSets {
			setMatched := true
			for _, pattern := range patternSet {
				if !strings.Contains(text, pattern) {
					setMatched = false
					break
				}
			}
			if setMatched {
				matched = true
				break
			}
		}
		if !matched {
			label := check.missingLabel
			if label == "" && len(check.patternSets) > 0 {
				label = strings.Join(check.patternSets[0], `", "`)
			}
			issues = append(issues, fmt.Sprintf("%s not clearly used; %s missing %q", check.label, check.file, label))
		}
	}

	if len(issues) > 0 {
		return failWithIssues("required stack", issues)
	}

	return pass("required stack", "required packages are declared and stack usage markers are present")
}
