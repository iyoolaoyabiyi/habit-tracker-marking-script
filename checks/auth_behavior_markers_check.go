package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkAuthBehaviorMarkers(root string) Result {
	authContent, err := os.ReadFile(filepath.Join(root, "src/lib/auth.ts"))
	if err != nil {
		return fail("auth behavior markers", err.Error())
	}
	loginContent, err := os.ReadFile(filepath.Join(root, "src/components/auth/LoginForm.tsx"))
	if err != nil {
		return fail("auth behavior markers", err.Error())
	}
	signupContent, err := os.ReadFile(filepath.Join(root, "src/components/auth/SignupForm.tsx"))
	if err != nil {
		return fail("auth behavior markers", err.Error())
	}
	dashboardContent, err := os.ReadFile(filepath.Join(root, "src/app/dashboard/page.tsx"))
	if err != nil {
		return fail("auth behavior markers", err.Error())
	}

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
		"src/lib/auth.ts":                    string(authContent),
		"src/components/auth/LoginForm.tsx":  string(loginContent),
		"src/components/auth/SignupForm.tsx": string(signupContent),
		"src/app/dashboard/page.tsx":         string(dashboardContent),
	}

	for _, item := range required {
		if !strings.Contains(sources[item.src], item.text) {
			return fail("auth behavior markers", fmt.Sprintf("%s missing %q", item.src, item.text))
		}
	}

	return pass("auth behavior markers", "auth error messages and redirects are clearly implemented")
}
