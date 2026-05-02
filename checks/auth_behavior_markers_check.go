package checks

import (
	"fmt"
)

func checkAuthBehaviorMarkers(root string) Result {
	var issues []string
	authContent, _ := readFileIssue(root, "src/lib/auth.ts", &issues)
	loginContent, _ := readFileIssue(root, "src/components/auth/LoginForm.tsx", &issues)
	signupContent, _ := readFileIssue(root, "src/components/auth/SignupForm.tsx", &issues)
	dashboardContent, _ := readFileIssue(root, "src/app/dashboard/page.tsx", &issues)

	required := []struct {
		text         string
		alternatives []string
		src          string
	}{
		{text: "User already exists", alternatives: []string{"User already exists"}, src: "src/lib/auth.ts"},
		{text: "Invalid email or password", alternatives: []string{"Invalid email or password"}, src: "src/lib/auth.ts"},
		{text: "/dashboard", alternatives: []string{"/dashboard", "DASHBOARD"}, src: "src/components/auth/LoginForm.tsx"},
		{text: "/dashboard", alternatives: []string{"/dashboard", "DASHBOARD"}, src: "src/components/auth/SignupForm.tsx"},
		{text: "/login", alternatives: []string{"/login", "LOGIN"}, src: "src/app/dashboard/page.tsx"},
	}

	sources := map[string]string{
		"src/lib/auth.ts":                    authContent,
		"src/components/auth/LoginForm.tsx":  loginContent,
		"src/components/auth/SignupForm.tsx": signupContent,
		"src/app/dashboard/page.tsx":         dashboardContent,
	}

	for _, item := range required {
		if !containsAny(sources[item.src], item.alternatives...) {
			issues = append(issues, fmt.Sprintf("%s missing %q", item.src, item.text))
		}
	}

	if len(issues) > 0 {
		return failWithIssues("auth behavior markers", issues)
	}

	return pass("auth behavior markers", "auth error messages and redirects are clearly implemented")
}
