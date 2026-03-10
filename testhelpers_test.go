package passgen

// testOpts returns a minimal opts for use in algorithm tests.
// It does not read environment variables.
func testOpts() opts {
	return opts{length: 40}
}
