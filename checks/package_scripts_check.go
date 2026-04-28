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
		return fail("package scripts", err.Error())
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fail("package scripts", "invalid package.json: "+err.Error())
	}

	required := map[string][]string{
		"dev":              {"next", "dev"},
		"build":            {"next", "build"},
		"start":            {"next", "start"},
		"test:unit":        {"vitest", "run", "coverage"},
		"test:integration": {"vitest", "run"},
		"test:e2e":         {"playwright", "test"},
		"test":             {"test:unit", "test:integration", "test:e2e"},
	}
	var missing, invalid []string
	for name, markers := range required {
		command, ok := pkg.Scripts[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		for _, marker := range markers {
			if !strings.Contains(command, marker) {
				invalid = append(invalid, fmt.Sprintf("%s missing %q", name, marker))
			}
		}
	}
	if len(missing) > 0 {
		return fail("package scripts", "missing script(s): "+strings.Join(missing, ", "))
	}
	if len(invalid) > 0 {
		return fail("package scripts", "script command marker mismatch: "+strings.Join(invalid, "; "))
	}

	return pass("package scripts", "all required script names and command markers present")
}
