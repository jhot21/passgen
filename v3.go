package main

import (
	"crypto/sha256"

	"golang.org/x/crypto/chacha20"
)

func passgenV3(opts opts) string {
	key := sha256.Sum256([]byte(opts.passphrase + opts.salt))
	nonce := make([]byte, chacha20.NonceSize) // 12 zero bytes

	stream, err := chacha20.NewUnauthenticatedCipher(key[:], nonce)
	if err != nil {
		// Cannot happen: key is always 32 bytes, nonce always 12 bytes
		panic(err)
	}

	classes := buildEnabledClasses(opts)
	if len(classes) == 0 {
		panic("passgen: all character classes are empty after filtering")
	}
	n := len(classes)
	result := make([]byte, opts.length)
	cycle := make([]string, n)

	for i := 0; i < opts.length; {
		copy(cycle, classes)
		shuffleClasses(stream, cycle)
		for j := 0; j < n && i < opts.length; j++ {
			result[i] = sampleFromClass(stream, cycle[j])
			i++
		}
	}
	return string(result)
}

func buildEnabledClasses(opts opts) []string {
	var classes []string
	ambiguous := ""
	if opts.noAmbiguous {
		ambiguous = defaultAmbiguous
		if opts.customAmbiguous != "" {
			ambiguous = opts.customAmbiguous
		}
	}

	maybeFilter := func(class string) string {
		if ambiguous == "" {
			return class
		}
		return filterAmbiguous(class, ambiguous)
	}

	if !opts.noUppers {
		if c := maybeFilter("ABCDEFGHIJKLMNOPQRSTUVWXYZ"); c != "" {
			classes = append(classes, c)
		}
	}
	if c := maybeFilter("abcdefghijklmnopqrstuvwxyz"); c != "" {
		classes = append(classes, c)
	}
	if !opts.noNumbers {
		if c := maybeFilter("0123456789"); c != "" {
			classes = append(classes, c)
		}
	}
	if !opts.noSpecials {
		specials := defaultSpecials
		if opts.customSpecials != "" {
			specials = opts.customSpecials
		}
		if c := maybeFilter(specials); c != "" {
			classes = append(classes, c)
		}
	}
	return classes
}

// sampleIndex returns a uniformly distributed integer in [0, n) using
// rejection sampling from stream, avoiding modulo bias.
func sampleIndex(stream *chacha20.Cipher, n int) int {
	limit := (256 / n) * n
	buf := make([]byte, 1)
	zeros := make([]byte, 1)
	for {
		stream.XORKeyStream(buf, zeros)
		b := int(buf[0])
		if b < limit {
			return b % n
		}
	}
}

// shuffleClasses performs a Fisher-Yates shuffle of classes in place,
// drawing swap indices from stream.
func shuffleClasses(stream *chacha20.Cipher, classes []string) {
	for i := len(classes) - 1; i > 0; i-- {
		j := sampleIndex(stream, i+1)
		classes[i], classes[j] = classes[j], classes[i]
	}
}

// sampleFromClass reads bytes from stream using rejection sampling to return
// a uniformly distributed character from class, avoiding modulo bias.
func sampleFromClass(stream *chacha20.Cipher, class string) byte {
	return class[sampleIndex(stream, len(class))]
}
