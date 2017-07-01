package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/davidlazar/go-libyaml/yaml"
	"github.com/davidlazar/pwclip"
	"golang.org/x/crypto/ssh/terminal"
)

func promptPassphrase() ([]byte, error) {
	fmt.Fprint(os.Stderr, "Passphrase: ")
	passphrase, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("terminal.ReadPassword: %s", err)
	}
	key, err := pwclip.Key(passphrase)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %s", err)
	}
	return key, nil
}

func newPWMFromYaml(yamlDoc []byte, question *int) (*pwclip.PWM, error) {
	pwm := &pwclip.PWM{
		Charset: pwclip.CharsetAlphaNumeric,
		Length:  32,
	}

	y, err := yaml.Load(yamlDoc)
	if err != nil {
		return nil, err
	}

	m, ok := y.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected yaml structure")
	}

	if v, ok := m["url"]; ok {
		if s, ok := v.(string); ok {
			pwm.URL = s
		} else {
			return nil, fmt.Errorf("url must be a string")
		}
	}
	if v, ok := m["username"]; ok {
		if s, ok := v.(string); ok {
			pwm.Username = s
		} else {
			return nil, fmt.Errorf("username must be a string")
		}
	}
	if v, ok := m["prefix"]; ok {
		if s, ok := v.(string); ok {
			pwm.Prefix = s
		} else {
			return nil, fmt.Errorf("prefix must be a string")
		}
	}
	if v, ok := m["charset"]; ok {
		if s, ok := v.(string); ok {
			pwm.Charset = s
		} else {
			return nil, fmt.Errorf("charset must be a string")
		}
	}
	if v, ok := m["length"]; ok {
		if s, ok := v.(string); ok {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("length must be an int")
			}
			pwm.Length = i
		} else {
			return nil, fmt.Errorf("length must be an int")
		}
	}
	if question != nil {
		q := "q" + strconv.Itoa(*question)
		v, ok := m[q]
		if !ok {
			return nil, fmt.Errorf("question %q not in settings", q)
		}
		if s, ok := v.(string); ok {
			pwm.Extra = &s
		} else {
			return nil, fmt.Errorf("q%d must be a string", *question)
		}
	}
	return pwm, nil
}
