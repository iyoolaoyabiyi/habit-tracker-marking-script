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
		return fail("required stack", err.Error())
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
	for _, name := range requiredPackages {
		if !hasPackage(name) {
			return fail("required stack", "package.json missing package "+name)
		}
	}

	stackUsageChecks := []struct {
		label    string
		file     string
		patterns []string
	}{
		{
			label:    "Next.js App Router",
			file:     "src/app/layout.tsx",
			patterns: []string{"export default function RootLayout"},
		},
		{
			label:    "React",
			file:     "src/components/auth/LoginForm.tsx",
			patterns: []string{"useState"},
		},
		{
			label:    "TypeScript",
			file:     "tsconfig.json",
			patterns: []string{`"compilerOptions"`},
		},
		{
			label:    "Tailwind CSS",
			file:     "src/app/globals.css",
			patterns: []string{`@import "tailwindcss"`},
		},
		{
			label:    "Playwright",
			file:     "playwright.config.ts",
			patterns: []string{"@playwright/test", "defineConfig"},
		},
		{
			label:    "Vitest",
			file:     "vitest.config.ts",
			patterns: []string{"vitest/config", "defineConfig"},
		},
		{
			label:    "React Testing Library",
			file:     "tests/integration/auth-flow.test.tsx",
			patterns: []string{"@testing-library/react"},
		},
	}

	for _, check := range stackUsageChecks {
		fileContent, err := os.ReadFile(filepath.Join(root, check.file))
		if err != nil {
			return fail("required stack", fmt.Sprintf("%s evidence missing: %v", check.label, err))
		}
		text := string(fileContent)
		for _, pattern := range check.patterns {
			if !strings.Contains(text, pattern) {
				return fail("required stack", fmt.Sprintf("%s not clearly used; %s missing %q", check.label, check.file, pattern))
			}
		}
	}

	return pass("required stack", "required packages are declared and stack usage markers are present")
}
