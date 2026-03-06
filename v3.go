package main

import (
	"crypto/sha256"

	"golang.org/x/crypto/chacha20"
)

func passgenV3(opts opts) string {
	// placeholder — implemented in Task 5
	return ""
}

func buildEnabledClasses(opts opts) []string {
	var classes []string
	if !opts.noUppers {
		classes = append(classes, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	classes = append(classes, "abcdefghijklmnopqrstuvwxyz")
	if !opts.noNumbers {
		classes = append(classes, "0123456789")
	}
	if !opts.noSpecials {
		specials := defaultSpecials
		if opts.customSpecials != "" {
			specials = opts.customSpecials
		}
		classes = append(classes, specials)
	}
	return classes
}

func sampleFromClass(stream *chacha20.Cipher, class string) byte {
	// placeholder — implemented in Task 5
	return 0
}

// suppress unused import errors until Task 5
var _ = sha256.Sum256
