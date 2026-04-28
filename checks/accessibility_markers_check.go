package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkAccessibilityMarkers(root string) Result {
	files := []string{
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
	}

	for _, rel := range files {
		content, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return fail("accessibility markers", err.Error())
		}
		text := string(content)
		if !strings.Contains(text, "<label") {
			return fail("accessibility markers", rel+" is missing label elements")
		}
		if !strings.Contains(text, "htmlFor=") {
			return fail("accessibility markers", rel+" is missing htmlFor associations")
		}
		if !strings.Contains(text, "<button") {
			return fail("accessibility markers", rel+" is missing button elements")
		}
	}

	globalCSS, err := os.ReadFile(filepath.Join(root, "src/app/globals.css"))
	if err != nil {
		return fail("accessibility markers", err.Error())
	}
	if !strings.Contains(string(globalCSS), ":focus-visible") {
		return fail("accessibility markers", "src/app/globals.css is missing visible focus styling")
	}

	return pass("accessibility markers", "labels, buttons, and visible focus styles are present")
}
