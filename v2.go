package main

import (
	"crypto/sha256"
	"strings"

	"github.com/eknkc/basex"
)

const defaultSpecials string = "!@#$%^&*"
const defaultAmbiguous string = "Il1O05S"

func passgenV2(opts opts) string {
	charset := createCharset(opts)
	encoding, _ := basex.NewEncoding(charset)
	var out string

	for i := 0; true; i++ {
		if len(out) >= opts.length { // password is long enough
			if containsNecessaryCharacters(out[0:opts.length], opts) { // password contains all necessary character classes
				break
			} else { // password is missing character classes, keep second half and keep generating if necessary
				out = out[int(len(out)/2):]
			}
		}

		passphrasePart := strings.Repeat(opts.passphrase, i+1)
		hash := sha256.Sum256([]byte(passphrasePart + opts.salt))
		out += encoding.Encode(hash[:])
	}

	return out[0:opts.length]
}

func createCharset(opts opts) string {
	charset := ""
	if !opts.noUppers {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	charset += "abcdefghijklmnopqrstuvwxyz"
	if !opts.noNumbers {
		charset += "0123456789"
	}
	if !opts.noSpecials {
		if opts.customSpecials != "" {
			charset += opts.customSpecials
		} else {
			charset += defaultSpecials
		}
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

func containsNecessaryCharacters(s string, opts opts) bool {
	if !opts.noNumbers && !strings.ContainsAny(s, "0123456789") {
		return false
	}

	if !opts.noUppers && !strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	if !opts.noSpecials {
		specials := defaultSpecials
		if opts.customSpecials != "" {
			specials = opts.customSpecials
		}
		if !strings.ContainsAny(s, specials) {
			return false
		}
	}

	return true
}
