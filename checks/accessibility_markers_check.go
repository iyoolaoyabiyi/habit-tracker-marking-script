package checks

import (
	"strings"
)

func checkAccessibilityMarkers(root string) Result {
	files := []string{
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
	}

	var issues []string
	for _, rel := range files {
		text, ok := readFileIssue(root, rel, &issues)
		if !ok {
			continue
		}
		if !strings.Contains(text, "<label") {
			issues = append(issues, rel+" is missing label elements")
		}
		if !strings.Contains(text, "htmlFor=") {
			issues = append(issues, rel+" is missing htmlFor associations")
		}
		if !strings.Contains(text, "<button") {
			issues = append(issues, rel+" is missing button elements")
		}
	}

	globalCSS, ok := readFileIssue(root, "src/app/globals.css", &issues)
	if ok && !strings.Contains(globalCSS, ":focus-visible") {
		issues = append(issues, "src/app/globals.css is missing :focus-visible styling")
	}

	if len(issues) > 0 {
		return failWithIssues("accessibility markers", issues)
	}

	return pass("accessibility markers", "labels, buttons, and visible focus styles are present")
}
