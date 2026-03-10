# passgen

A deterministic password generator: the same passphrase + salt always produces the same password. Instead of storing passwords, you remember simple passphrases and derive site-specific passwords on demand.

## Install

```sh
go install codeberg.org/jhot/passgen/cmd/passgen@latest
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
    -v, --version           algorithm version (default: 3)
    -s, --salt              salt, overrides $PASSGEN_SALT
    -l, --length            output length (default: 40, or $PASSGEN_LENGTH)
    -p, --print             print to terminal instead of copying to clipboard
        --no-uppers         exclude uppercase letters
        --no-numbers        exclude digits
        --no-specials       exclude special characters
        --custom-specials   use a custom set of special characters
        --no-ambiguous      exclude visually ambiguous characters (Il1O05S)
        --custom-ambiguous  use a custom set of ambiguous characters to exclude
```

By default the generated password is copied to the clipboard using the [OSC 52](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands) terminal escape sequence. This works over SSH and regardless of display server (X11, Wayland), but requires a supported terminal emulator:

- **Supported:** kitty, WezTerm, Ghostty, foot, Alacritty, Windows Terminal, iTerm2 (enable in Preferences → General → "Applications in terminal may access clipboard"), xterm (requires `-allow` flag or `XTerm.allowWindowOps: true`)
- **tmux:** add `set -g set-clipboard on` to `~/.tmux.conf`
- **Not supported:** macOS Terminal.app, GNOME Terminal, and other VTE-based terminals (Terminator, XFCE Terminal, MATE Terminal)

*Note this list is not exhaustive or updated regularly, see https://can-i-use-terminal.github.io/features/osc52copy.html for a more up to date list*

If your terminal is not supported, use `--print` and pipe the output instead:

```sh
passgen --print github | pbcopy    # macOS
passgen --print github | xclip -selection clipboard  # Linux X11
```

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

### v3 (default)

`passphrase + salt` → SHA-256 → 32-byte ChaCha20 key → unlimited deterministic byte stream → output generated in cycles of N characters (where N = number of enabled classes). At the start of each cycle, the class order is shuffled via Fisher-Yates using bytes from the ChaCha20 stream, then one character is sampled from each class using rejection sampling for uniform distribution. Guarantees each enabled class (uppercase, lowercase, digits, specials) appears exactly once per cycle with no modulo bias and no predictable positional pattern.

### v2 (legacy)

`passphrase + salt` → SHA-256 → basex-encoded using the active character set. If the output is too short or missing a required character class, the process repeats with an extended passphrase and the results are concatenated. Use `--version 2` to select it.

### v1 (legacy)

`passphrase + salt` → UTF-16 LE → Skein-512 → hex → base91-encoded → substring at a hash-derived offset. Compatible with [node-passgen-cli](https://github.com/jhotmann/node-passgen-cli). Use `--version 1` to select it.

## Environment variables

| Variable | Description |
|---|---|
| `PASSGEN_SALT` | Default salt appended to every passphrase |
| `PASSGEN_LENGTH` | Default output length (overrides the built-in default of 40) |
