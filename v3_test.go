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

func TestPassgenV3_Deterministic(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	p1 := passgenV3(o)
	p2 := passgenV3(o)
	if p1 != p2 {
		t.Errorf("passgenV3 is not deterministic: %q != %q", p1, p2)
	}
}

func TestPassgenV3_Length(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	p := passgenV3(o)
	if len(p) != 40 {
		t.Errorf("expected length 40, got %d", len(p))
	}
}

func TestPassgenV3_LongPassword(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 200
	p := passgenV3(o)
	if len(p) != 200 {
		t.Errorf("expected length 200, got %d", len(p))
	}
}

func TestPassgenV3_ContainsAllClasses(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	p := passgenV3(o)

	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false
	for _, c := range p {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	if !hasUpper {
		t.Error("password missing uppercase")
	}
	if !hasLower {
		t.Error("password missing lowercase")
	}
	if !hasDigit {
		t.Error("password missing digit")
	}
	if !hasSpecial {
		t.Error("password missing special")
	}
}

func TestPassgenV3_NoUppers(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	o.noUppers = true
	p := passgenV3(o)
	for _, c := range p {
		if c >= 'A' && c <= 'Z' {
			t.Errorf("password contains uppercase %q but noUppers=true", c)
		}
	}
}

func TestPassgenV3_NoNumbers(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	o.noNumbers = true
	p := passgenV3(o)
	for _, c := range p {
		if c >= '0' && c <= '9' {
			t.Errorf("password contains digit %q but noNumbers=true", c)
		}
	}
}

func TestPassgenV3_NoSpecials(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	o.noSpecials = true
	p := passgenV3(o)
	for _, c := range p {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			t.Errorf("password contains special %q but noSpecials=true", c)
		}
	}
}

func TestPassgenV3_ProportionalDistribution(t *testing.T) {
	// Generate a long password and verify each class appears roughly equally.
	// With 4 classes and 400 chars, expect ~100 per class. Allow ±30 tolerance.
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 400
	p := passgenV3(o)

	counts := map[string]int{"upper": 0, "lower": 0, "digit": 0, "special": 0}
	for _, c := range p {
		switch {
		case c >= 'A' && c <= 'Z':
			counts["upper"]++
		case c >= 'a' && c <= 'z':
			counts["lower"]++
		case c >= '0' && c <= '9':
			counts["digit"]++
		default:
			counts["special"]++
		}
	}
	expected := 100
	tolerance := 30
	for class, count := range counts {
		diff := count - expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("class %q count %d is outside tolerance (expected ~%d ±%d)", class, count, expected, tolerance)
		}
	}
}

func TestPassgenV3_DifferentInputsDifferentOutputs(t *testing.T) {
	o1 := defaultOpts()
	o1.passphrase = "passphrase1"
	o1.salt = "salt"

	o2 := defaultOpts()
	o2.passphrase = "passphrase2"
	o2.salt = "salt"

	if passgenV3(o1) == passgenV3(o2) {
		t.Error("different passphrases should produce different passwords")
	}
}
