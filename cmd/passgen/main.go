package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	passgen "codeberg.org/jhot/passgen"
	osc52 "github.com/aymanbagabas/go-osc52/v2"
	"github.com/peterh/liner"
	"github.com/teris-io/cli"
)

func main() {
	app := cli.New("generate secure passwords from a simple passphrase").
		WithArg(cli.NewArg("passphrase", "input passphrase").AsOptional()).
		WithOption(cli.NewOption("version", "override passgen algorithm version").WithChar('v').WithType(cli.TypeInt)).
		WithOption(cli.NewOption("salt", "salt appended to passphrase, default env[PASSGEN_SALT]").WithChar('s').WithType(cli.TypeString)).
		WithOption(cli.NewOption("length", "password length, default 40 or env[PASSGEN_LENGTH]").WithChar('l').WithType(cli.TypeInt)).
		WithOption(cli.NewOption("print", "print password instead of copying").WithChar('p').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("no-specials", "no special characters").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("no-uppers", "no uppercase characters").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("no-numbers", "no number characters").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("custom-specials", "custom special character set").WithType(cli.TypeString)).
		WithOption(cli.NewOption("no-ambiguous", "exclude ambiguous characters (Il1O05S)").WithType(cli.TypeBool)).
		WithOption(cli.NewOption("custom-ambiguous", "custom set of ambiguous characters to exclude").WithType(cli.TypeString)).
		WithAction(func(args []string, options map[string]string) int {
			var passphrase string
			if len(args) > 0 {
				passphrase = strings.Join(args, " ")
			} else {
				passphrase = promptForPassphrase()
			}
			if passphrase == "" {
				return 1
			}

			b := passgen.New(passphrase)

			// Salt: flag takes priority, then env var
			if saltValue, ok := options["salt"]; ok {
				b = b.WithSalt(saltValue)
			} else if envSalt := os.Getenv("PASSGEN_SALT"); envSalt != "" {
				b = b.WithSalt(envSalt)
			}

			// Length: flag takes priority, then env var
			if lengthValue, ok := options["length"]; ok {
				if n, err := strconv.Atoi(lengthValue); err == nil {
					b = b.WithLength(n)
				}
			} else if envLen := os.Getenv("PASSGEN_LENGTH"); envLen != "" {
				if n, err := strconv.Atoi(envLen); err == nil && n > 0 {
					b = b.WithLength(n)
				}
			}

			if versionValue, ok := options["version"]; ok {
				if v, err := strconv.Atoi(versionValue); err == nil {
					b = b.WithVersion(v)
				}
			}

			if customSpecials, ok := options["custom-specials"]; ok {
				b = b.WithSpecials(customSpecials)
			}

			if customAmbiguous, ok := options["custom-ambiguous"]; ok {
				b = b.WithAmbiguousChars(customAmbiguous)
			} else if _, ok := options["no-ambiguous"]; ok {
				b = b.NoAmbiguous()
			}

			if _, ok := options["no-uppers"]; ok {
				b = b.NoUppers()
			}
			if _, ok := options["no-numbers"]; ok {
				b = b.NoNumbers()
			}
			if _, ok := options["no-specials"]; ok {
				b = b.NoSpecials()
			}

			pass, err := b.Generate()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				return 1
			}

			if _, ok := options["print"]; ok {
				fmt.Println(pass)
			} else {
				osc52.New(pass).WriteTo(os.Stderr)
			}
			return 0
		})

	os.Exit(app.Run(os.Args, os.Stdout))
}

func promptForPassphrase() string {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()

	if passphrase, err := line.PasswordPrompt("Passphrase: "); err == nil {
		return passphrase
	}
	return ""
}
