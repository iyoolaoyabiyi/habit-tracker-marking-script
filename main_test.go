package main

import "testing"

func TestMatchGitHubTreeBranchUsesLongestBranchPrefix(t *testing.T) {
	refs := "abc\trefs/heads/main\n" +
		"def\trefs/heads/feature/stage-3\n" +
		"ghi\trefs/heads/feature\n"

	branch, subdir := matchGitHubTreeBranch(refs, "feature/stage-3/frontend")

	if branch != "feature/stage-3" {
		t.Fatalf("expected longest branch match, got %q", branch)
	}
	if subdir != "frontend" {
		t.Fatalf("expected subdir frontend, got %q", subdir)
	}
}

func TestMatchGitHubTreeBranchHandlesBranchRoot(t *testing.T) {
	refs := "abc\trefs/heads/main\n"

	branch, subdir := matchGitHubTreeBranch(refs, "main")

	if branch != "main" {
		t.Fatalf("expected branch main, got %q", branch)
	}
	if subdir != "" {
		t.Fatalf("expected empty subdir, got %q", subdir)
	}
}
