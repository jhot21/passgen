package passgen_test

import (
	"strings"
	"testing"

	"codeberg.org/jhot/passgen"
)

// Happy path: basic generation with v3 (default)
func TestBuilder_Generate_HappyPath(t *testing.T) {
	pass, err := passgen.New("mypassphrase").WithSalt("mysalt").WithLength(20).Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pass) != 20 {
		t.Errorf("expected length 20, got %d", len(pass))
	}
}

// Determinism: same inputs always produce the same output
func TestBuilder_Generate_Deterministic(t *testing.T) {
	a, _ := passgen.New("phrase").WithSalt("salt").Generate()
	b, _ := passgen.New("phrase").WithSalt("salt").Generate()
	if a != b {
		t.Errorf("expected deterministic output, got %q and %q", a, b)
	}
}

// Different inputs produce different outputs
func TestBuilder_Generate_DifferentSalts(t *testing.T) {
	a, _ := passgen.New("phrase").WithSalt("salt1").Generate()
	b, _ := passgen.New("phrase").WithSalt("salt2").Generate()
	if a == b {
		t.Errorf("expected different outputs for different salts")
	}
}

// Error: empty passphrase
func TestBuilder_Generate_EmptyPassphrase(t *testing.T) {
	_, err := passgen.New("").Generate()
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

// Error: length < 1
func TestBuilder_Generate_ZeroLength(t *testing.T) {
	_, err := passgen.New("phrase").WithLength(0).Generate()
	if err == nil {
		t.Fatal("expected error for zero length")
	}
}

// Error: invalid version
func TestBuilder_Generate_InvalidVersion(t *testing.T) {
	_, err := passgen.New("phrase").WithVersion(99).Generate()
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}

// NoUppers+NoNumbers+NoSpecials produces lowercase-only password, no error
func TestBuilder_Generate_LowercaseOnly(t *testing.T) {
	pass, err := passgen.New("phrase").NoUppers().NoNumbers().NoSpecials().Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pass) != 40 {
		t.Errorf("expected length 40, got %d", len(pass))
	}
	for _, c := range pass {
		if c < 'a' || c > 'z' {
			t.Errorf("expected lowercase only, got char %q in password %q", c, pass)
		}
	}
}

// Error: all character classes truly empty (all lowercase filtered as ambiguous + all others disabled)
func TestBuilder_Generate_AllClassesEmpty(t *testing.T) {
	_, err := passgen.New("phrase").
		WithAmbiguousChars("abcdefghijklmnopqrstuvwxyz").
		NoUppers().NoNumbers().NoSpecials().
		Generate()
	if err == nil {
		t.Fatal("expected error when all character classes are empty after filtering")
	}
}

// Builder immutability: chaining from the same base doesn't mutate it
func TestBuilder_Immutability(t *testing.T) {
	base := passgen.New("phrase").WithSalt("salt")
	a, _ := base.WithLength(20).Generate()
	b, _ := base.WithLength(30).Generate()
	if len(a) != 20 {
		t.Errorf("expected length 20, got %d", len(a))
	}
	if len(b) != 30 {
		t.Errorf("expected length 30, got %d", len(b))
	}
	// original base should still have default length (40)
	c, _ := base.Generate()
	if len(c) != 40 {
		t.Errorf("expected base length 40, got %d", len(c))
	}
}

// Version 1 works via builder
func TestBuilder_Generate_V1(t *testing.T) {
	pass, err := passgen.New("phrase").WithSalt("salt").WithVersion(1).Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pass) != 40 {
		t.Errorf("expected length 40, got %d", len(pass))
	}
}

// Version 2 works via builder
func TestBuilder_Generate_V2(t *testing.T) {
	pass, err := passgen.New("phrase").WithSalt("salt").WithVersion(2).Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pass) != 40 {
		t.Errorf("expected length 40, got %d", len(pass))
	}
}

// WithSpecials sets a custom special character set
func TestBuilder_WithSpecials(t *testing.T) {
	pass, err := passgen.New("phrase").WithSpecials("!").Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.ContainsAny(pass, "!") {
		t.Errorf("expected password to contain '!', got %q", pass)
	}
}

// NoAmbiguous removes default ambiguous chars
func TestBuilder_NoAmbiguous(t *testing.T) {
	const ambiguous = "Il1O05S"
	for range 20 {
		pass, err := passgen.New("phrase").WithSalt("salt").NoAmbiguous().Generate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.ContainsAny(pass, ambiguous) {
			t.Errorf("password contains ambiguous chars: %q", pass)
		}
	}
}

// WithAmbiguousChars uses a custom ambiguous set
func TestBuilder_WithAmbiguousChars(t *testing.T) {
	for range 20 {
		pass, err := passgen.New("phrase").WithSalt("salt").WithAmbiguousChars("xyz").Generate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.ContainsAny(pass, "xyz") {
			t.Errorf("password contains custom ambiguous chars: %q", pass)
		}
	}
}
