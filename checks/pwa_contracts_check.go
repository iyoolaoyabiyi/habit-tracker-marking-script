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
		return fail("pwa contract", err.Error())
	}

	var mf manifest
	if err := json.Unmarshal(content, &mf); err != nil {
		return fail("pwa contract", "invalid manifest.json: "+err.Error())
	}

	if mf.Name == "" || mf.ShortName == "" || mf.StartURL == "" || mf.Display == "" || mf.BackgroundColor == "" || mf.ThemeColor == "" {
		return fail("pwa contract", "manifest.json is missing one or more required fields")
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
		return fail("pwa contract", "manifest.json is missing required 192/512 icons")
	}

	swContent, err := os.ReadFile(filepath.Join(root, "public/sw.js"))
	if err != nil {
		return fail("pwa contract", err.Error())
	}
	if !strings.Contains(string(swContent), "caches") || !strings.Contains(string(swContent), "fetch") {
		return fail("pwa contract", "service worker does not appear to cache app shell requests")
	}

	registrationText, err := readExistingFiles(root, []string{
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/components/shared/ProtectedRoute.tsx",
		"src/components/shared/SplashScreen.tsx",
	})
	if err != nil {
		return fail("pwa contract", err.Error())
	}
	if !strings.Contains(registrationText, "serviceWorker") || !strings.Contains(registrationText, "sw.js") {
		return fail("pwa contract", "service worker registration not found in required app/shared files")
	}

	return pass("pwa contract", "manifest, icons, service worker, and registration found")
}
