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
		return skip("executable verification", "skipped because -run-commands=false")
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
	coverageOutput, err := verifyCoverageSummary(root)
	if err != nil {
		return fail("executable verification", "coverage threshold verification failed:\n"+tail(coverageOutput, 60))
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

	runtimeOutput, err := runExaminerRuntimeTests(root)
	if err != nil {
		return fail("executable verification", "examiner-owned runtime tests failed:\n"+tail(runtimeOutput, 80))
	}

	var browserOutput string
	if options.RunE2E && (options.RunOffline || options.RunAccessibility || options.RunResponsive) {
		buildOutput, err := runCommand(root, "npm", "run", "build")
		if err != nil {
			return fail("executable verification", "npm run build before examiner browser tests failed:\n"+tail(buildOutput, 40))
		}

		browserOutput, err = runExaminerBrowserTests(root, options)
		if err != nil {
			return fail("executable verification", "examiner-owned browser tests failed:\n"+tail(browserOutput, 100))
		}
	}

	details := "build, unit, integration, and examiner runtime tests ran successfully"
	if options.RunE2E {
		if options.RunOffline || options.RunAccessibility || options.RunResponsive {
			details += "; e2e and selected examiner browser tests ran successfully"
		} else {
			details += "; e2e ran successfully; examiner browser checks skipped by flags"
		}
	} else {
		details += "; e2e and examiner browser tests skipped because -run-e2e=false"
	}

	return pass("executable verification", details)
}
