package checks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkPackageScripts(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return failWithIssues("package scripts", []string{"package.json could not be read: " + err.Error()})
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fail("package scripts", "invalid package.json: "+err.Error())
	}

	required := []struct {
		name    string
		markers []string
	}{
		{name: "dev", markers: []string{"next", "dev"}},
		{name: "build", markers: []string{"next", "build"}},
		{name: "start", markers: []string{"next", "start"}},
		{name: "test:unit", markers: []string{"vitest", "run", "coverage"}},
		{name: "test:integration", markers: []string{"vitest", "run"}},
		{name: "test:e2e", markers: []string{"playwright", "test"}},
		{name: "test", markers: []string{"test:unit", "test:integration", "test:e2e"}},
	}
	var missing, invalid []string
	for _, check := range required {
		command, ok := pkg.Scripts[check.name]
		if !ok {
			missing = append(missing, check.name)
			continue
		}
		for _, marker := range check.markers {
			if !strings.Contains(command, marker) {
				invalid = append(invalid, fmt.Sprintf("%s missing %q", check.name, marker))
			}
		}
	}
	if len(missing) > 0 {
		for _, name := range missing {
			invalid = append(invalid, "package.json missing script "+name)
		}
	}
	if len(invalid) > 0 {
		return failWithIssues("package scripts", invalid)
	}

	return pass("package scripts", "all required script names and command markers present")
}
