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

func containsNecessaryCharacters(s string, o opts) bool {
	if !o.noNumbers && !strings.ContainsAny(s, "0123456789") {
		return false
	}

	if !o.noUppers && !strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	if !o.noSpecials {
		specials := defaultSpecials
		if o.customSpecials != "" {
			specials = o.customSpecials
		}
		if !strings.ContainsAny(s, specials) {
			return false
		}
	}

	return true
}
