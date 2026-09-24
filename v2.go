package passgen

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/eknkc/basex"
)

const defaultSpecials string = "!@#$%^&*"

// I (upper I), l (lower L), 1 (one), O (upper O), 0 (zero), 5, S
const defaultAmbiguous string = "Il1O05S"

func passgenV2(o opts) (string, error) {
	charset := createCharset(o)
	encoding, err := basex.NewEncoding(charset)
	if err != nil {
		return "", fmt.Errorf("passgen: invalid charset: %w", err)
	}
	if err := checkV2Satisfiable(o, charset); err != nil {
		return "", err
	}
	var out string

	for i := 0; true; i++ {
		if len(out) >= o.length { // password is long enough
			if containsNecessaryCharacters(out[0:o.length], o) { // password contains all necessary character classes
				break
			} else { // password is missing character classes, keep second half and keep generating if necessary
				out = out[int(len(out)/2):]
			}
		}

		passphrasePart := strings.Repeat(o.passphrase, i+1)
		hash := sha256.Sum256([]byte(passphrasePart + o.salt))
		out += encoding.Encode(hash[:])
	}

	return out[0:o.length], nil
}

func createCharset(o opts) string {
	charset := ""
	if !o.noUppers {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	charset += "abcdefghijklmnopqrstuvwxyz"
	if !o.noNumbers {
		charset += "0123456789"
	}
	if !o.noSpecials {
		if o.customSpecials != "" {
			charset += o.customSpecials
		} else {
			charset += defaultSpecials
		}
	}

	if o.noAmbiguous {
		ambiguous := defaultAmbiguous
		if o.customAmbiguous != "" {
			ambiguous = o.customAmbiguous
		}
		charset = filterAmbiguous(charset, ambiguous)
	}

	return charset
}

// filterAmbiguous returns class with any characters present in ambiguous removed.
func filterAmbiguous(class, ambiguous string) string {
	result := make([]byte, 0, len(class))
	for i := 0; i < len(class); i++ {
		if !strings.ContainsRune(ambiguous, rune(class[i])) {
			result = append(result, class[i])
		}
	}
	return string(result)
}

// requiredClassesV2 returns the character classes a v2 password must contain,
// unfiltered, exactly as containsNecessaryCharacters checks them.
func requiredClassesV2(o opts) []string {
	var classes []string
	if !o.noNumbers {
		classes = append(classes, "0123456789")
	}
	if !o.noUppers {
		classes = append(classes, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if !o.noSpecials {
		specials := defaultSpecials
		if o.customSpecials != "" {
			specials = o.customSpecials
		}
		classes = append(classes, specials)
	}
	return classes
}

// checkV2Satisfiable returns an error when the passgenV2 loop could never
// terminate: an alphabet too small to encode with, a required class removed
// entirely by ambiguous filtering, or a length too short to hold one character
// from every required class.
func checkV2Satisfiable(o opts, charset string) error {
	if len([]rune(charset)) < 2 {
		return fmt.Errorf("passgen: character set must contain at least 2 characters")
	}
	required := requiredClassesV2(o)
	for _, class := range required {
		if !strings.ContainsAny(charset, class) {
			return fmt.Errorf("passgen: a required character class is empty after ambiguous filtering")
		}
	}
	if o.length < len(required) {
		return fmt.Errorf("passgen: length %d is too short for %d required character classes", o.length, len(required))
	}
	return nil
}

func containsNecessaryCharacters(s string, o opts) bool {
	for _, class := range requiredClassesV2(o) {
		if !strings.ContainsAny(s, class) {
			return false
		}
	}
	return true
}
