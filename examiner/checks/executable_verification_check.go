package checks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkExecutableVerification(root string, options RuntimeOptions) Result {
	if !options.RunCommands {
		return pass("executable verification", "skipped because -run-commands=false")
	}

	if _, err := exec.LookPath("npm"); err != nil {
		return fail("executable verification", "npm is required to run build and test scripts")
	}

	if options.InstallDeps {
		installCmd := []string{"npm", "ci"}
		if _, err := os.Stat(filepath.Join(root, "package-lock.json")); err != nil {
			installCmd = []string{"npm", "install"}
		}
		output, err := runCommand(root, installCmd[0], installCmd[1:]...)
		if err != nil {
			return fail("executable verification", fmt.Sprintf("%s failed:\n%s", strings.Join(installCmd, " "), tail(output, 40)))
		}
	}

	buildOutput, err := runCommand(root, "npm", "run", "build")
	if err != nil {
		return fail("executable verification", "npm run build failed:\n"+tail(buildOutput, 40))
	}

	unitOutput, err := runCommand(root, "npm", "run", "test:unit")
	if err != nil {
		return fail("executable verification", "npm run test:unit failed:\n"+tail(unitOutput, 60))
	}
	if !strings.Contains(unitOutput, "Coverage") && !strings.Contains(unitOutput, "coverage") {
		return fail("executable verification", "unit test output is missing coverage output")
	}

	integrationOutput, err := runCommand(root, "npm", "run", "test:integration")
	if err != nil {
		return fail("executable verification", "npm run test:integration failed:\n"+tail(integrationOutput, 60))
	}

	if options.RunE2E {
		e2eOutput, err := runCommand(root, "npm", "run", "test:e2e")
		if err != nil {
			return fail("executable verification", "npm run test:e2e failed:\n"+tail(e2eOutput, 80))
		}
	}

	details := "build, unit, and integration scripts ran successfully"
	if options.RunE2E {
		details += "; e2e script ran successfully"
	} else {
		details += "; e2e script skipped because -run-e2e=false"
	}

	return pass("executable verification", details)
}
