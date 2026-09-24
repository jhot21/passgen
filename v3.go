package passgen

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/chacha20"
)

func passgenV3(o opts) (string, error) {
	key := sha256.Sum256([]byte(o.passphrase + o.salt))
	nonce := make([]byte, chacha20.NonceSize)

	stream, err := chacha20.NewUnauthenticatedCipher(key[:], nonce)
	if err != nil {
		return "", fmt.Errorf("passgen: cipher init failed: %w", err)
	}

	classes := buildEnabledClasses(o)
	if len(classes) == 0 {
		return "", fmt.Errorf("passgen: all character classes are empty after filtering")
	}
	for _, c := range classes {
		if len(c) > 256 {
			return "", fmt.Errorf("passgen: character class longer than 256 characters")
		}
	}
	n := len(classes)
	result := make([]byte, o.length)
	cycle := make([]string, n)

	for i := 0; i < o.length; {
		copy(cycle, classes)
		shuffleClasses(stream, cycle)
		for j := 0; j < n && i < o.length; j++ {
			result[i] = sampleFromClass(stream, cycle[j])
			i++
		}
	}
	return string(result), nil
}

func buildEnabledClasses(o opts) []string {
	var classes []string
	ambiguous := ""
	if o.noAmbiguous {
		ambiguous = defaultAmbiguous
		if o.customAmbiguous != "" {
			ambiguous = o.customAmbiguous
		}
	}

	maybeFilter := func(class string) string {
		if ambiguous == "" {
			return class
		}
		return filterAmbiguous(class, ambiguous)
	}

	if !o.noUppers {
		if c := maybeFilter("ABCDEFGHIJKLMNOPQRSTUVWXYZ"); c != "" {
			classes = append(classes, c)
		}
	}
	if c := maybeFilter("abcdefghijklmnopqrstuvwxyz"); c != "" {
		classes = append(classes, c)
	}
	if !o.noNumbers {
		if c := maybeFilter("0123456789"); c != "" {
			classes = append(classes, c)
		}
	}
	if !o.noSpecials {
		specials := defaultSpecials
		if o.customSpecials != "" {
			specials = o.customSpecials
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
