package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	osc52 "github.com/aymanbagabas/go-osc52/v2"
	"github.com/peterh/liner"
	"github.com/teris-io/cli"
)

func main() {
	// Define CLI options
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
			// Parse options
			opts := defaultOpts()

			versionValue, versionExists := options["version"]
			if versionExists {
				parsed, err := strconv.ParseInt(versionValue, 10, 0)
				if err == nil {
					opts.version = int(parsed)
				}
			}

			saltValue, saltExists := options["salt"]
			if saltExists {
				opts.salt = saltValue
			}

			customSpecialsValue, customSpecialsExists := options["custom-specials"]
			if customSpecialsExists {
				opts.customSpecials = customSpecialsValue
			}

			lengthValue, lengthExists := options["length"]
			if lengthExists {
				parsed, err := strconv.ParseInt(lengthValue, 10, 0)
				if err == nil {
					opts.length = int(parsed)
				}
			} else {
				parsed, err := strconv.ParseInt(os.Getenv("PASSGEN_LENGTH"), 10, 0)
				if err == nil && parsed > 0 {
					opts.length = int(parsed)
				}
			}

			_, opts.noSpecials = options["no-specials"]
			_, opts.noUppers = options["no-uppers"]
			_, opts.noNumbers = options["no-numbers"]
			_, opts.noAmbiguous = options["no-ambiguous"]

			customAmbiguousValue, customAmbiguousExists := options["custom-ambiguous"]
			if customAmbiguousExists {
				opts.customAmbiguous = customAmbiguousValue
				opts.noAmbiguous = true
			}

			if len(args) > 0 { // passphrase supplied as an argument
				opts.passphrase = strings.Join(args, " ")
			} else { // prompt for passphrase
				opts.passphrase = promptForPassword()
			}

			if opts.passphrase == "" {
				return 1
			}

			// Generate password using chosen algorithm
			pass := generatePassword(opts)
			// Print or copy password
			_, printPass := options["print"]
			if printPass {
				fmt.Println(pass)
			} else {
				osc52.New(pass).WriteTo(os.Stderr)
			}
			return 0
		})

	os.Exit(app.Run(os.Args, os.Stdout))
}

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

func defaultOpts() opts {
	return opts{
		salt:   os.Getenv("PASSGEN_SALT"),
		length: 40,
	}
}

func generatePassword(opts opts) string {
	switch opts.version {
	case 1:
		return passgenV1(opts)
	case 2:
		return passgenV2(opts)
	default:
		return passgenV3(opts)
	}
}

func promptForPassword() string {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer line.Close()

	if passphrase, err := line.PasswordPrompt("Passphrase: "); err == nil {
		return passphrase
	} else {
		return ""
	}
}
