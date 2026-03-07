package main

import (
	"crypto/sha256"
	"testing"

	"golang.org/x/crypto/chacha20"
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
	if len(p1) == 0 {
		t.Fatal("passgenV3 returned empty string")
	}
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
	if len(p) == 0 {
		t.Fatal("passgenV3 returned empty string")
	}
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
	if len(p) == 0 {
		t.Fatal("passgenV3 returned empty string")
	}
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
	if len(p) == 0 {
		t.Fatal("passgenV3 returned empty string")
	}
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

func TestPassgenV3_GoldenValue(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	o.length = 40
	got := passgenV3(o)
	// Golden value — regenerate with: go run . --salt=mysalt --print "correct horse battery staple"
	// If this changes, the v3 algorithm has changed and existing passwords will differ.
	want := "Ap4@Cq1#Fe9*Qv5#Wi8@Bk8*Pc0@Fw9@Nf1*Qs3&"
	if got != want {
		t.Errorf("v3 golden value changed\ngot:  %q\nwant: %q", got, want)
	}
}

func TestSampleIndex_InRange(t *testing.T) {
	key := sha256.Sum256([]byte("test"))
	nonce := make([]byte, chacha20.NonceSize)
	stream, _ := chacha20.NewUnauthenticatedCipher(key[:], nonce)
	for i := 0; i < 1000; i++ {
		idx := sampleIndex(stream, 4)
		if idx < 0 || idx >= 4 {
			t.Fatalf("sampleIndex(%d) = %d, want [0,4)", i, idx)
		}
	}
}

func TestSampleIndex_UniformDistribution(t *testing.T) {
	key := sha256.Sum256([]byte("test"))
	nonce := make([]byte, chacha20.NonceSize)
	stream, _ := chacha20.NewUnauthenticatedCipher(key[:], nonce)
	counts := make([]int, 4)
	const trials = 4000
	for i := 0; i < trials; i++ {
		counts[sampleIndex(stream, 4)]++
	}
	expected := trials / 4
	tolerance := expected / 5
	for i, c := range counts {
		diff := c - expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("index %d: count %d outside tolerance (expected ~%d ±%d)", i, c, expected, tolerance)
		}
	}
}

func TestShuffleClasses_IsPermutation(t *testing.T) {
	key := sha256.Sum256([]byte("test"))
	nonce := make([]byte, chacha20.NonceSize)
	stream, _ := chacha20.NewUnauthenticatedCipher(key[:], nonce)

	original := []string{
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789",
		defaultSpecials,
	}
	classes := make([]string, len(original))
	copy(classes, original)

	shuffleClasses(stream, classes)

	found := make(map[string]bool)
	for _, c := range classes {
		found[c] = true
	}
	for _, c := range original {
		if !found[c] {
			t.Errorf("element %q missing after shuffle", c)
		}
	}

	// Verify the shuffle actually changed the order (with this fixed seed it must)
	allSame := true
	for i := range classes {
		if classes[i] != original[i] {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("shuffleClasses did not change the order of elements")
	}
}

func TestShuffleClasses_Deterministic(t *testing.T) {
	// Two identical streams must produce identical shuffle orders.
	newStream := func() *chacha20.Cipher {
		key := sha256.Sum256([]byte("test"))
		nonce := make([]byte, chacha20.NonceSize)
		s, _ := chacha20.NewUnauthenticatedCipher(key[:], nonce)
		return s
	}

	input := []string{
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789",
		defaultSpecials,
	}

	a := make([]string, len(input))
	b := make([]string, len(input))
	copy(a, input)
	copy(b, input)

	shuffleClasses(newStream(), a)
	shuffleClasses(newStream(), b)

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("position %d: got %q and %q; same seed must produce same order", i, a[i], b[i])
		}
	}
}

func TestPassgenV3_CycleOrderVaries(t *testing.T) {
	o := defaultOpts()
	o.passphrase = "correct horse battery staple"
	o.salt = "mysalt"
	n := len(buildEnabledClasses(o))
	o.length = 2 * n // exactly two full cycles

	p := passgenV3(o)
	if len(p) < 2*n {
		t.Fatalf("expected password length >= %d, got %d", 2*n, len(p))
	}

	classOf := func(c rune) string {
		switch {
		case c >= 'A' && c <= 'Z':
			return "U"
		case c >= 'a' && c <= 'z':
			return "L"
		case c >= '0' && c <= '9':
			return "D"
		default:
			return "S"
		}
	}

	runes := []rune(p)
	cycle1 := make([]string, n)
	cycle2 := make([]string, n)
	for i := range n {
		cycle1[i] = classOf(runes[i])
		cycle2[i] = classOf(runes[i+n])
	}

	t.Logf("cycle1: %v", cycle1)
	t.Logf("cycle2: %v", cycle2)

	for i := range cycle1 {
		if cycle1[i] != cycle2[i] {
			return // cycles differ — test passes
		}
	}
	t.Errorf("cycle1 and cycle2 have identical class order %v; expected variation", cycle1)
}
