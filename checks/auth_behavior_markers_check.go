package checks

import (
	"fmt"
	"strings"
)

func checkAuthBehaviorMarkers(root string) Result {
	var issues []string
	authContent, _ := readFileIssue(root, "src/lib/auth.ts", &issues)
	loginContent, _ := readFileIssue(root, "src/components/auth/LoginForm.tsx", &issues)
	signupContent, _ := readFileIssue(root, "src/components/auth/SignupForm.tsx", &issues)
	dashboardContent, _ := readFileIssue(root, "src/app/dashboard/page.tsx", &issues)

	required := []struct {
		text string
		src  string
	}{
		{"User already exists", "src/lib/auth.ts"},
		{"Invalid email or password", "src/lib/auth.ts"},
		{"/dashboard", "src/components/auth/LoginForm.tsx"},
		{"/dashboard", "src/components/auth/SignupForm.tsx"},
		{"/login", "src/app/dashboard/page.tsx"},
	}

	sources := map[string]string{
		"src/lib/auth.ts":                    authContent,
		"src/components/auth/LoginForm.tsx":  loginContent,
		"src/components/auth/SignupForm.tsx": signupContent,
		"src/app/dashboard/page.tsx":         dashboardContent,
	}

	for _, item := range required {
		if !strings.Contains(sources[item.src], item.text) {
			issues = append(issues, fmt.Sprintf("%s missing %q", item.src, item.text))
		}
	}

	if len(issues) > 0 {
		return failWithIssues("auth behavior markers", issues)
	}

	return pass("auth behavior markers", "auth error messages and redirects are clearly implemented")
}
