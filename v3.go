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
	result := make([]byte, opts.length)
	for i := range result {
		class := classes[i%len(classes)]
		result[i] = sampleFromClass(stream, class)
	}
	return string(result)
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

// sampleFromClass reads bytes from stream using rejection sampling to return
// a uniformly distributed index into class, avoiding modulo bias.
func sampleFromClass(stream *chacha20.Cipher, class string) byte {
	size := len(class)
	// largest multiple of size that fits in a byte (0-255)
	limit := (256 / size) * size
	buf := make([]byte, 1)
	zeros := make([]byte, 1)
	for {
		stream.XORKeyStream(buf, zeros)
		b := int(buf[0])
		if b < limit {
			return class[b%size]
		}
		// byte is in the bias zone — discard and try again
	}
}
