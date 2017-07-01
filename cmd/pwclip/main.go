package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	printPass = false // -p
	question  *int    // -q <int>
	keyPath   *string // -k <file>
	yamlPath  *string // <file>
)

const version = "4.0"
const shortUsage = "usage: %s [-k <keyfile>] [-q <num>] [-p] <yamlfile>\n"
const longUsage = `
Required argument:
  <yamlfile>  password settings in YAML format

Options:
  -h --help
    Show this help message and exit.

  -v --version
    Print version number and exit.

  -k <keyfile>
    Read key from file, instead of prompting for a passphrase.

  -p
    Print password to stdout, instead of copying it to the clipboard.

  -q <num>
    Produce answer to the selected secret question.

More information: https://github.com/davidlazar/pwclip
`

func main() {
	log.SetFlags(0)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")

	parseArgs(os.Args[1:])

	yamlDoc, err := ioutil.ReadFile(*yamlPath)
	check(err, "file")

	pwm, err := newPWMFromYaml(yamlDoc, question)
	check(err, "settings")

	var key []byte
	if keyPath != nil {
		key, err = ioutil.ReadFile(*keyPath)
	} else {
		key, err = promptPassphrase()
	}
	check(err, "read key")

	pw := pwm.Password(key)

	if printPass {
		fmt.Println(pw)
	} else {
		fmt.Fprintf(os.Stderr, "Password copied to clipboard for 10 seconds.\n")
		check(setClipboardTemporarily([]byte(pw), 10*time.Second), "clipboard")
	}
}

func parseArgs(args []string) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			if yamlPath == nil {
				yamlPath = &s
				continue
			} else {
				urr("unrecognized argument: %q", s)
			}
		}
		switch {
		case s == "-h" || s == "--help":
			fmt.Printf(shortUsage, os.Args[0])
			fmt.Print(longUsage)
			os.Exit(0)
		case s == "-v" || s == "--version":
			fmt.Println(version)
			os.Exit(0)
		case s == "-p":
			printPass = true
		case strings.HasPrefix(s, "-k"):
			k := parseArg("-k", s, &args)
			keyPath = &k
		case strings.HasPrefix(s, "-q"):
			q := parseArg("-q", s, &args)
			if qi, err := strconv.Atoi(q); err != nil {
				urr("flag -q: invalid number value: %q", q)
			} else {
				question = &qi
			}
		default:
			urr("unrecognized flag: %q", s)
		}
	}
	if yamlPath == nil {
		urr("missing required argument: <yamlfile>")
	}
}

func parseArg(flag string, arg string, args *[]string) string {
	if flag == arg {
		if len(*args) == 0 {
			urr("flag %s: expected argument", flag)
		}
		r := (*args)[0]
		*args = (*args)[1:]
		return r
	}
	return arg[len(flag):]
}

func check(err error, prefix string) {
	if err != nil {
		log.Fatalf("error (%s): %s", prefix, err)
	}
}

// usage error
func urr(format string, a ...interface{}) {
	log.Printf("error (usage): "+format, a...)
	fmt.Printf(shortUsage, os.Args[0])
	os.Exit(1)
}
