// Package passgen generates deterministic passwords from a passphrase and salt.
// The same passphrase + salt + options always produces the same password.
//
// Example:
//
//	pass, err := passgen.New("mypassphrase").WithSalt("mysalt").WithLength(32).Generate()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(pass)
package passgen

import "fmt"

// Builder constructs a deterministic password generation request.
// Use [New] to create a Builder with sensible defaults.
//
// Builder methods return a new Builder (copy), so a base Builder can be
// reused safely without mutation:
//
//	base := passgen.New("mypassphrase").WithSalt("mysalt")
//	pass16, _ := base.WithLength(16).Generate()
//	pass32, _ := base.WithLength(32).Generate()
type Builder struct {
	o opts
}

// New returns a Builder for the given passphrase with defaults:
// v3 algorithm, length 40, all character classes enabled,
// specials = !@#$%^&*. No environment variables are read.
func New(passphrase string) *Builder {
	return &Builder{o: opts{
		passphrase: passphrase,
		length:     40,
	}}
}

// WithVersion selects the algorithm version. Valid values: 1, 2, 3.
// The default (0 or unset) maps to version 3.
func (b *Builder) WithVersion(v int) *Builder {
	c := b.copy()
	c.o.version = v
	return c
}

// WithSalt appends a salt to the passphrase before hashing.
// A unique salt makes passwords specific to a user or context.
func (b *Builder) WithSalt(salt string) *Builder {
	c := b.copy()
	c.o.salt = salt
	return c
}

// WithLength sets the output password length. Must be at least 1. Default: 40.
func (b *Builder) WithLength(n int) *Builder {
	c := b.copy()
	c.o.length = n
	return c
}

// WithSpecials sets a custom special character set to use instead of the
// default (!@#$%^&*). Has no effect if [NoSpecials] has been called.
func (b *Builder) WithSpecials(chars string) *Builder {
	c := b.copy()
	c.o.customSpecials = chars
	return c
}

// WithAmbiguousChars sets a custom set of characters to exclude as visually
// ambiguous. Calling this implicitly enables ambiguous filtering.
// The default ambiguous set (used by [NoAmbiguous]) is "Il1O05S".
func (b *Builder) WithAmbiguousChars(chars string) *Builder {
	c := b.copy()
	c.o.customAmbiguous = chars
	c.o.noAmbiguous = true
	return c
}

// NoUppers disables uppercase letters in the generated password.
func (b *Builder) NoUppers() *Builder {
	c := b.copy()
	c.o.noUppers = true
	return c
}

// NoNumbers disables numeric characters in the generated password.
func (b *Builder) NoNumbers() *Builder {
	c := b.copy()
	c.o.noNumbers = true
	return c
}

// NoSpecials disables special characters in the generated password entirely.
func (b *Builder) NoSpecials() *Builder {
	c := b.copy()
	c.o.noSpecials = true
	return c
}

// NoAmbiguous excludes the default ambiguous characters (Il1O05S) from the
// generated password. Use [WithAmbiguousChars] to supply a custom set.
func (b *Builder) NoAmbiguous() *Builder {
	c := b.copy()
	c.o.noAmbiguous = true
	return c
}

// Generate produces a deterministic password from the passphrase provided to [New].
//
// Returns an error if:
//   - The passphrase is empty.
//   - The length is less than 1.
//   - The version is not 1, 2, or 3.
//   - All character classes are disabled (leaving no characters to sample from).
//   - (v2) The character set has fewer than 2 characters.
//   - (v2) A required character class is filtered away entirely (e.g. by ambiguous-character exclusion).
//   - (v2) The length is shorter than the number of required character classes.
//   - (v3) A character class is longer than 256 characters.
func (b *Builder) Generate() (string, error) {
	if b.o.passphrase == "" {
		return "", fmt.Errorf("passgen: passphrase must not be empty")
	}
	if b.o.length < 1 {
		return "", fmt.Errorf("passgen: length must be at least 1")
	}
	switch b.o.version {
	case 0, 3: // 0 is the Go zero value — defaults to v3
		return passgenV3(b.o)
	case 1:
		return passgenV1(b.o)
	case 2:
		return passgenV2(b.o)
	default:
		return "", fmt.Errorf("passgen: unsupported version %d (valid: 1, 2, 3)", b.o.version)
	}
}

func (b *Builder) copy() *Builder {
	c := *b
	return &c
}

// opts holds the internal configuration for password generation.
// It is unexported; callers use Builder methods to configure it.
type opts struct {
	version         int
	passphrase      string
	salt            string
	length          int
	noUppers        bool
	noNumbers       bool
	noSpecials      bool
	customSpecials  string
	noAmbiguous     bool
	customAmbiguous string
}
