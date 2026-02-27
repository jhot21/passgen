# passgen

A deterministic password generator: the same passphrase + salt always produces the same password. Instead of storing passwords, you remember simple passphrases and derive site-specific passwords on demand.

## Install

```sh
go install codeberg.org/jhot/passgen@latest
```

## Setup

Set a personal salt. This makes your passwords unique to you — even a trivial passphrase like `github` produces a completely different result than it would for anyone else.

Generate one with passgen itself:

```sh
passgen --no-specials -l 80 somerandomstring
# E9Fz0Fxs1JTND8E6HxJvcLjMJK6ZaH2DjzY8Pmot85mtLwX8Z1QaCCo6sUWrHlP7gavv9aBs3MQ9WPQ1
```

Then set it in your shell profile (`~/.zshrc`, `~/.bashrc`, etc.):

```sh
export PASSGEN_SALT="your-long-random-string-here"
```

On Windows, set `PASSGEN_SALT` in System → Advanced System Settings → Environment Variables.

## Usage

```
passgen [passphrase]

Arguments:
    passphrase              input passphrase (prompted if omitted)

Options:
    -v, --version           algorithm version (default: 2)
    -s, --salt              salt, overrides $PASSGEN_SALT
    -l, --length            output length (default: 40, or $PASSGEN_LENGTH)
    -p, --print             print to terminal instead of copying to clipboard
        --no-uppers         exclude uppercase letters
        --no-numbers        exclude digits
        --no-specials       exclude special characters
        --custom-specials   use a custom set of special characters
```

By default the generated password is copied to the clipboard silently. Use `--print` to display it instead.

**Examples:**

```sh
# Generate a password for GitHub (copied to clipboard)
passgen github

# Print instead of copy
passgen --print github

# Site that only allows alphanumeric passwords
passgen --no-specials github

# Site with a restricted special character set
passgen --custom-specials '!@#$' github

# Shorter password
passgen --length 16 github
```

## Algorithms

### v2 (default)

`passphrase + salt` → SHA-256 → basex-encoded using the active character set. If the output is too short or missing a required character class, the process repeats with an extended passphrase and the results are concatenated. Always produces uppercase, lowercase, digits, and specials unless disabled.

### v1 (legacy)

`passphrase + salt` → UTF-16 LE → Skein-512 → hex → base91-encoded → substring at a hash-derived offset. Compatible with [node-passgen-cli](https://github.com/jhotmann/node-passgen-cli). Use `--version 1` to select it.

## Environment variables

| Variable | Description |
|---|---|
| `PASSGEN_SALT` | Default salt appended to every passphrase |
| `PASSGEN_LENGTH` | Default output length (overrides the built-in default of 40) |
