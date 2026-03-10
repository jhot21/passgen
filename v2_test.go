package passgen

import (
	"strings"
	"testing"
)

// --- passgenV2 regression tests ---

func TestPassgenV2Regression(t *testing.T) {
	tests := []struct {
		name string
		o    opts
		want string
	}{
		{
			name: "default options",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40},
			want: "OMTZbN0GP9!d3KvA0^oPZ2#nlML2yc%*bEw!3!Hd",
		},
		{
			name: "empty salt",
			o:    opts{passphrase: "mypassphrase", salt: "", length: 40},
			want: "MQETO$fgy9OmGEBU5Rxp5s40UuHivbv@x3X&62&w",
		},
		{
			name: "no specials",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40, noSpecials: true},
			want: "hHfG2KBdWnZjDXOxoUMGxElH04DaDhWQM1HTrRxs",
		},
		{
			name: "no uppers",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40, noUppers: true},
			want: "p%gd*yeh@1t&c4fjsius^3x!6y$5r&cddp4%3z#n",
		},
		{
			name: "no numbers",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40, noNumbers: true},
			want: "CLQzhhC&Hl&VWoRzhObkQVLHj$xgKsRkXEblB$oj",
		},
		{
			name: "all classes disabled",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40, noUppers: true, noNumbers: true, noSpecials: true},
			want: "cmdamckvagqgfvurtrxleapgsdfpxpvwouowjlxk",
		},
		{
			name: "length 20",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 20},
			want: "OMTZbN0GP9!d3KvA0^oP",
		},
		{
			name: "custom specials",
			o:    opts{passphrase: "mypassphrase", salt: "testsalt", length: 40, customSpecials: "+-="},
			want: "yFB586PXgb6jEsWoIwDG4LBYtStN=0Y4p2hI7j9g",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := passgenV2(tt.o)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- createCharset ---

func TestCreateCharset(t *testing.T) {
	uppers := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowers := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	specials := "!@#$%^&*"

	tests := []struct {
		name string
		o    opts
		want string
	}{
		{"default", opts{}, uppers + lowers + digits + specials},
		{"no uppers", opts{noUppers: true}, lowers + digits + specials},
		{"no numbers", opts{noNumbers: true}, uppers + lowers + specials},
		{"no specials", opts{noSpecials: true}, uppers + lowers + digits},
		{"all disabled", opts{noUppers: true, noNumbers: true, noSpecials: true}, lowers},
		{"custom specials", opts{customSpecials: "+-="}, uppers + lowers + digits + "+-="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createCharset(tt.o)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- containsNecessaryCharacters ---

func TestContainsNecessaryCharacters(t *testing.T) {
	tests := []struct {
		name string
		s    string
		o    opts
		want bool
	}{
		// default (requires uppers, numbers, specials)
		{"all classes present", "Password1!", opts{}, true},
		{"missing uppercase", "password1!", opts{}, false},
		{"missing number", "Password!", opts{}, false},
		{"missing special", "Password1", opts{}, false},
		// noUppers — requires numbers and specials only
		{"noUppers all met", "password1!", opts{noUppers: true}, true},
		{"noUppers missing number", "password!", opts{noUppers: true}, false},
		{"noUppers missing special", "password1", opts{noUppers: true}, false},
		// noNumbers — requires uppers and specials only
		{"noNumbers all met", "Password!", opts{noNumbers: true}, true},
		{"noNumbers missing upper", "password!", opts{noNumbers: true}, false},
		// noSpecials — requires uppers and numbers only
		{"noSpecials all met", "Password1", opts{noSpecials: true}, true},
		{"noSpecials missing upper", "password1", opts{noSpecials: true}, false},
		// all disabled — always true
		{"all disabled accepts anything", "alllowercase", opts{noUppers: true, noNumbers: true, noSpecials: true}, true},
		// custom specials
		{"custom specials present", "Password1+", opts{customSpecials: "+-="}, true},
		{"custom specials missing", "Password1!", opts{customSpecials: "+-="}, false},
		// regex-special chars in custom specials set
		{"custom specials with hyphen", "Password1-", opts{customSpecials: "-"}, true},
		{"custom specials with caret", "Password1^", opts{customSpecials: "^"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsNecessaryCharacters(tt.s, tt.o)
			if got != tt.want {
				t.Errorf("containsNecessaryCharacters(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestFilterAmbiguous_RemovesChars(t *testing.T) {
	got := filterAmbiguous("ABCDEFGHIJKLMNOPQRSTUVWXYZ", defaultAmbiguous)
	for _, c := range got {
		if strings.ContainsRune(defaultAmbiguous, c) {
			t.Errorf("filterAmbiguous left ambiguous char %q in result %q", c, got)
		}
	}
	if strings.Contains(got, "O") {
		t.Errorf("expected O removed, got %q", got)
	}
	if !strings.Contains(got, "A") {
		t.Errorf("expected A present, got %q", got)
	}
}

func TestFilterAmbiguous_EmptyResult(t *testing.T) {
	got := filterAmbiguous("Il1", defaultAmbiguous)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFilterAmbiguous_LowercaseAlphabet(t *testing.T) {
	got := filterAmbiguous("abcdefghijklmnopqrstuvwxyz", defaultAmbiguous)
	// defaultAmbiguous contains lowercase l — it should be removed
	if strings.Contains(got, "l") {
		t.Errorf("expected l removed, got %q", got)
	}
	if !strings.Contains(got, "a") {
		t.Errorf("expected a present, got %q", got)
	}
}

func TestCreateCharset_NoAmbiguous(t *testing.T) {
	o := testOpts()
	o.passphrase = "test"
	o.noAmbiguous = true
	charset := createCharset(o)
	for _, c := range charset {
		if strings.ContainsRune(defaultAmbiguous, c) {
			t.Errorf("charset contains ambiguous char %q: %q", c, charset)
		}
	}
}

func TestCreateCharset_CustomAmbiguous(t *testing.T) {
	o := testOpts()
	o.passphrase = "test"
	o.noAmbiguous = true
	o.customAmbiguous = "abc"
	charset := createCharset(o)
	// Custom set replaces default: abc should be absent
	for _, c := range charset {
		if strings.ContainsRune("abc", c) {
			t.Errorf("charset contains custom ambiguous char %q: %q", c, charset)
		}
	}
	// defaultAmbiguous chars should still be present (not filtered)
	// 'l' is in defaultAmbiguous but not in customAmbiguous="abc"
	if !strings.ContainsRune(charset, 'l') {
		t.Errorf("expected 'l' (defaultAmbiguous char) to be present when using customAmbiguous, got %q", charset)
	}
}
