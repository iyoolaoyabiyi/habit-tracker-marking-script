package checks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func checkPWAContracts(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "public/manifest.json"))
	if err != nil {
		return failWithIssues("pwa contract", []string{"public/manifest.json could not be read: " + err.Error()})
	}

	var mf manifest
	if err := json.Unmarshal(content, &mf); err != nil {
		return failWithIssues("pwa contract", []string{"public/manifest.json is invalid JSON: " + err.Error()})
	}

	var issues []string
	if mf.Name == "" {
		issues = append(issues, "manifest.json is missing name")
	}
	if mf.ShortName == "" {
		issues = append(issues, "manifest.json is missing short_name")
	}
	if mf.StartURL == "" {
		issues = append(issues, "manifest.json is missing start_url")
	}
	if mf.Display == "" {
		issues = append(issues, "manifest.json is missing display")
	}
	if mf.BackgroundColor == "" {
		issues = append(issues, "manifest.json is missing background_color")
	}
	if mf.ThemeColor == "" {
		issues = append(issues, "manifest.json is missing theme_color")
	}

	var has192, has512 bool
	for _, icon := range mf.Icons {
		if strings.Contains(icon.Src, "icon-192.png") && strings.Contains(icon.Sizes, "192x192") {
			has192 = true
		}
		if strings.Contains(icon.Src, "icon-512.png") && strings.Contains(icon.Sizes, "512x512") {
			has512 = true
		}
	}
	if !has192 || !has512 {
		if !has192 {
			issues = append(issues, "manifest.json is missing required 192x192 icon")
		}
		if !has512 {
			issues = append(issues, "manifest.json is missing required 512x512 icon")
		}
	}

	swContent, err := os.ReadFile(filepath.Join(root, "public/sw.js"))
	if err != nil {
		issues = append(issues, "public/sw.js could not be read: "+err.Error())
	} else {
		swText := string(swContent)
		if !strings.Contains(swText, "caches") {
			issues = append(issues, "public/sw.js is missing caches usage")
		}
		if !strings.Contains(swText, "fetch") {
			issues = append(issues, "public/sw.js is missing fetch handling")
		}
	}

	registrationText, err := readExistingFiles(root, []string{
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/components/shared/ProtectedRoute.tsx",
		"src/components/shared/SplashScreen.tsx",
	})
	if err != nil {
		issues = append(issues, "required app/shared files for service worker registration could not be read: "+err.Error())
	} else {
		if !strings.Contains(registrationText, "serviceWorker") {
			issues = append(issues, "service worker registration missing serviceWorker marker")
		}
		if !strings.Contains(registrationText, "sw.js") {
			issues = append(issues, "service worker registration missing sw.js marker")
		}
	}

	if len(issues) > 0 {
		return failWithIssues("pwa contract", issues)
	}

	return pass("pwa contract", "manifest, icons, service worker, and registration found")
}
