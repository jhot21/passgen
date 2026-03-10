# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build the CLI
go build ./cmd/passgen/

# Install CLI locally
go install ./cmd/passgen/

# Install CLI from anywhere
# go install codeberg.org/jhot/passgen/cmd/passgen@latest

# Run CLI directly (for development)
go run ./cmd/passgen/ [passphrase]

# Run tests
go test ./...

# Update dependencies
go mod tidy
```

## Architecture

`passgen` is a deterministic password generator: the same passphrase + salt always produces the same password. The tool either copies the result to the clipboard or prints it.

**File layout:**
- [builder.go](builder.go) — public `Builder` type, all exported methods, and the unexported `opts` struct; this is the library's public API
- [v1.go](v1.go) — v1 algorithm: passphrase+salt → UTF-16 LE → Skein-512 hash → hex → base91 encoded; uses a hash-derived start offset to pick the substring
- [v2.go](v2.go) — v2 algorithm: passphrase+salt → SHA-256 → basex-encoded using a custom charset; iterates with repeated passphrase until the output is long enough and satisfies all required character classes (uppercase, numbers, specials)
- [v3.go](v3.go) — v3 algorithm (default): ChaCha20 keystream derived from SHA-256(passphrase+salt); rejection sampling with Fisher-Yates shuffle of character classes per cycle
- [cmd/passgen/main.go](cmd/passgen/main.go) — CLI entry point; reads `PASSGEN_SALT`/`PASSGEN_LENGTH` env vars; uses Builder API

**Algorithm selection:** `Generate()` in [builder.go](builder.go) dispatches on the configured version; anything other than `1` or `2` uses v3.

**Environment variables:**
- `PASSGEN_SALT` — default salt (strongly recommended to set; makes passwords unique to you)
- `PASSGEN_LENGTH` — default output length (overrides the built-in default of 40)
