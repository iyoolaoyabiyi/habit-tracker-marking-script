package checks

import (
	"fmt"
	"strings"
)

type enforcementKind string

const (
	enforcementStatic    enforcementKind = "static"
	enforcementRuntime   enforcementKind = "runtime"
	enforcementBrowser   enforcementKind = "browser"
	enforcementHeuristic enforcementKind = "heuristic"
)

type specRequirement struct {
	Section     string
	Description string
	Kind        enforcementKind
	CheckGroups []string
}

var specChecklist = []specRequirement{
	{Section: "3", Description: "required stack and no remote backend", Kind: enforcementStatic, CheckGroups: []string{"required stack", "no remote backend markers"}},
	{Section: "4", Description: "route files and core redirects", Kind: enforcementStatic, CheckGroups: []string{"route contract files", "splash timing", "source behavior markers", "executable verification"}},
	{Section: "5", Description: "localStorage keys, persistence access, and storage ownership markers", Kind: enforcementStatic, CheckGroups: []string{"storage keys", "local persistence usage", "source behavior markers"}},
	{Section: "6", Description: "required folder and file structure", Kind: enforcementStatic, CheckGroups: []string{"required files"}},
	{Section: "7", Description: "naming conventions", Kind: enforcementStatic, CheckGroups: []string{"naming conventions"}},
	{Section: "8", Description: "exported auth and habit type contracts", Kind: enforcementStatic, CheckGroups: []string{"type contracts"}},
	{Section: "9", Description: "utility function contracts", Kind: enforcementStatic, CheckGroups: []string{"utility contracts", "executable verification"}},
	{Section: "10", Description: "UI contract test ids", Kind: enforcementStatic, CheckGroups: []string{"ui contract markers"}},
	{Section: "11", Description: "auth behavior markers", Kind: enforcementStatic, CheckGroups: []string{"auth behavior markers", "source behavior markers", "executable verification"}},
	{Section: "12", Description: "habit behavior markers", Kind: enforcementStatic, CheckGroups: []string{"habit behavior markers", "source behavior markers", "executable verification"}},
	{Section: "13", Description: "PWA manifest, service worker, and registration", Kind: enforcementStatic, CheckGroups: []string{"pwa contract", "executable verification"}},
	{Section: "14", Description: "responsive usability", Kind: enforcementBrowser, CheckGroups: []string{"executable verification"}},
	{Section: "15", Description: "accessibility markers", Kind: enforcementHeuristic, CheckGroups: []string{"accessibility markers"}},
	{Section: "16", Description: "required test files and exact test titles", Kind: enforcementStatic, CheckGroups: []string{"test suite contract", "executable verification"}},
	{Section: "17", Description: "coverage configuration and coverage output", Kind: enforcementRuntime, CheckGroups: []string{"coverage config", "executable verification"}},
	{Section: "18", Description: "required package scripts", Kind: enforcementStatic, CheckGroups: []string{"package scripts", "executable verification"}},
	{Section: "19", Description: "README sections", Kind: enforcementStatic, CheckGroups: []string{"readme requirements"}},
}

func checkSpecChecklist() Result {
	if len(specChecklist) == 0 {
		return fail("spec checklist", "no requirements are registered")
	}

	seenSections := map[string]bool{}
	kindCounts := map[enforcementKind]int{}
	var missing []string
	for _, item := range specChecklist {
		seenSections[item.Section] = true
		kindCounts[item.Kind]++
		if item.Description == "" || len(item.CheckGroups) == 0 {
			missing = append(missing, "section "+item.Section)
		}
	}
	if len(missing) > 0 {
		return fail("spec checklist", "requirements without descriptions or check groups: "+strings.Join(missing, ", "))
	}

	for _, section := range []string{"3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19"} {
		if !seenSections[section] {
			return fail("spec checklist", "missing requirement coverage for section "+section)
		}
	}

	return pass(
		"spec checklist",
		fmt.Sprintf(
			"%d requirements classified: static=%d, runtime=%d, browser=%d, heuristic=%d",
			len(specChecklist),
			kindCounts[enforcementStatic],
			kindCounts[enforcementRuntime],
			kindCounts[enforcementBrowser],
			kindCounts[enforcementHeuristic],
		),
	)
}
