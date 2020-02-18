package main

import "testing"

func TestRemoveFeatureFromBranch(t *testing.T) {
	f := RemoveFeatureFromBranch("feature/001-feature-name")
	if f != "001-feature-name" {
		t.Errorf("remove feature incorrect, got %s, want: %s ", f, "001-feature-name")
	}
}

func TestVersionToA(t *testing.T) {
	ver:= Version{Major: 1,Minor: 2, Patch: 3}
	f := VersionToA(&ver)
	if f != "1.2.3" {
		t.Errorf("version conversion error, got %s, want: %s", f, "1.2.3")
	}
}
