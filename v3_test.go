package main

import (
	"testing"
)

func TestBuildEnabledClasses_Default(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	classes := buildEnabledClasses(o)
	if len(classes) != 4 {
		t.Fatalf("expected 4 classes, got %d", len(classes))
	}
	// Order: uppercase, lowercase, digits, specials
	if classes[0] != "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		t.Errorf("expected uppercase class, got %q", classes[0])
	}
	if classes[1] != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("expected lowercase class, got %q", classes[1])
	}
	if classes[2] != "0123456789" {
		t.Errorf("expected digits class, got %q", classes[2])
	}
	if classes[3] != defaultSpecials {
		t.Errorf("expected specials class, got %q", classes[3])
	}
}

func TestBuildEnabledClasses_NoUppers(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	o.noUppers = true
	classes := buildEnabledClasses(o)
	if len(classes) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(classes))
	}
	if classes[0] != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("expected lowercase first, got %q", classes[0])
	}
}

func TestBuildEnabledClasses_NoNumbers(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	o.noNumbers = true
	classes := buildEnabledClasses(o)
	if len(classes) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(classes))
	}
	for _, c := range classes {
		if c == "0123456789" {
			t.Error("digits class should not be present when noNumbers=true")
		}
	}
}

func TestBuildEnabledClasses_NoSpecials(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	o.noSpecials = true
	classes := buildEnabledClasses(o)
	if len(classes) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(classes))
	}
	for _, c := range classes {
		if c == defaultSpecials {
			t.Error("specials class should not be present when noSpecials=true")
		}
	}
}

func TestBuildEnabledClasses_OnlyLowercase(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	o.noUppers = true
	o.noNumbers = true
	o.noSpecials = true
	classes := buildEnabledClasses(o)
	if len(classes) != 1 {
		t.Fatalf("expected 1 class, got %d", len(classes))
	}
	if classes[0] != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("expected lowercase, got %q", classes[0])
	}
}

func TestBuildEnabledClasses_CustomSpecials(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "test"
	o.customSpecials = "+-="
	classes := buildEnabledClasses(o)
	if len(classes) != 4 {
		t.Fatalf("expected 4 classes, got %d", len(classes))
	}
	if classes[3] != "+-=" {
		t.Errorf("expected custom specials, got %q", classes[3])
	}
}
