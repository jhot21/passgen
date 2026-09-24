// Command goldengen prints golden test vectors for ports of passgen v2/v3.
// Usage: go run ./internal/goldengen > golden_vectors.json
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jhot21/passgen"
)

type vector struct {
	Name            string `json:"name"`
	Version         int    `json:"version"`
	Passphrase      string `json:"passphrase"`
	Salt            string `json:"salt"`
	Length          int    `json:"length"`
	NoUppers        bool   `json:"noUppers"`
	NoNumbers       bool   `json:"noNumbers"`
	NoSpecials      bool   `json:"noSpecials"`
	CustomSpecials  string `json:"customSpecials"`
	NoAmbiguous     bool   `json:"noAmbiguous"`
	CustomAmbiguous string `json:"customAmbiguous"`
	Expected        string `json:"expected"`
	Error           bool   `json:"error"`
}

func run(v vector) vector {
	b := passgen.New(v.Passphrase).WithVersion(v.Version).WithSalt(v.Salt).WithLength(v.Length)
	if v.NoUppers {
		b = b.NoUppers()
	}
	if v.NoNumbers {
		b = b.NoNumbers()
	}
	if v.NoSpecials {
		b = b.NoSpecials()
	}
	if v.CustomSpecials != "" {
		b = b.WithSpecials(v.CustomSpecials)
	}
	if v.NoAmbiguous {
		b = b.NoAmbiguous()
	}
	if v.CustomAmbiguous != "" {
		b = b.WithAmbiguousChars(v.CustomAmbiguous)
		v.NoAmbiguous = true
	}
	out, err := b.Generate()
	if err != nil {
		v.Error = true
		v.Expected = ""
	} else {
		v.Expected = out
	}
	return v
}

func main() {
	inputs := []struct{ pass, salt string }{
		{"correct horse", ""},
		{"hunter2", "salt"},
		{"pässwörd✓", "sälz"},
	}
	customs := []struct{ specials, ambiguous string }{
		{"", ""},
		{"-_.", ""},
		{"", "abc"},
		{"!!", ""},  // duplicate: v2 errors, v3 fine
		{"abc", ""}, // overlaps lowercase: v2 errors, v3 fine
	}
	var out []vector
	for _, version := range []int{2, 3} {
		for mask := 0; mask < 16; mask++ {
			for _, length := range []int{5, 14, 40, 128} {
				for _, in := range inputs {
					for _, c := range customs {
						v := vector{
							Version: version, Passphrase: in.pass, Salt: in.salt, Length: length,
							NoUppers: mask&1 != 0, NoNumbers: mask&2 != 0,
							NoSpecials: mask&4 != 0, NoAmbiguous: mask&8 != 0,
							CustomSpecials: c.specials, CustomAmbiguous: c.ambiguous,
						}
						v.Name = fmt.Sprintf("v%d mask%02d len%d %q/%q specials=%q ambiguous=%q",
							version, mask, length, in.pass, in.salt, c.specials, c.ambiguous)
						out = append(out, run(v))
					}
				}
			}
		}
	}
	errorCases := []vector{
		{Name: "err v2 specials filtered away", Version: 2, Passphrase: "p", Length: 20, CustomSpecials: "l", CustomAmbiguous: "l"},
		{Name: "err v2 length below required", Version: 2, Passphrase: "p", Length: 2},
		{Name: "err v2 one-char alphabet", Version: 2, Passphrase: "p", Length: 10, NoUppers: true, NoNumbers: true, NoSpecials: true, CustomAmbiguous: "bcdefghijklmnopqrstuvwxyz"},
		{Name: "err v3 class over 256", Version: 3, Passphrase: "p", Length: 10, CustomSpecials: strings.Repeat("!", 257)},
		{Name: "err v3 all filtered", Version: 3, Passphrase: "p", Length: 10, NoUppers: true, NoNumbers: true, NoSpecials: true, CustomAmbiguous: "abcdefghijklmnopqrstuvwxyz"},
	}
	for _, v := range errorCases {
		out = append(out, run(v))
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
