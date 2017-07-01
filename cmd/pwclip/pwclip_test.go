package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/davidlazar/pwclip"
)

var settingsTests = []struct {
	question *int
	yaml     string
	pwm      *pwclip.PWM
	err      error
}{
	{
		nil, `
url: example.com
username: example@example.com
q1: frequent flier number
q2: first car
odd: {ignored: true}`,
		&pwclip.PWM{
			URL:      "example.com",
			Username: "example@example.com",
			Charset:  pwclip.CharsetAlphaNumeric,
			Length:   32,
		},
		nil,
	},
	{
		nil, `# comment
url: example.com # comment
username: example@example.com
q1: frequent flier number
# comment
q2: first car`,
		&pwclip.PWM{
			URL:      "example.com",
			Username: "example@example.com",
			Charset:  pwclip.CharsetAlphaNumeric,
			Length:   32,
		},
		nil,
	},
	{
		nil, `
url:
username: 
prefix: # comment`,
		&pwclip.PWM{
			URL:      "",
			Username: "",
			Prefix:   "",
			Charset:  pwclip.CharsetAlphaNumeric,
			Length:   32,
		},
		nil,
	},
	{
		nil, `{url: example.com, username: example@example.com, q1: frequent flier number, q2: first car}`,
		&pwclip.PWM{
			URL:      "example.com",
			Username: "example@example.com",
			Charset:  pwclip.CharsetAlphaNumeric,
			Length:   32,
		},
		nil,
	},
	{
		intptr(2), `
url: example.com
username: example@example.com
q1: frequent flier number
q2: first car`,
		&pwclip.PWM{
			URL:      "example.com",
			Username: "example@example.com",
			Extra:    strptr("first car"),
			Charset:  pwclip.CharsetAlphaNumeric,
			Length:   32,
		},
		nil,
	},
	{
		nil, `
username: example@example.com
charset: αβγδεζηθικλμνξοπρστυφχψω`,
		&pwclip.PWM{
			Username: "example@example.com",
			Charset:  "αβγδεζηθικλμνξοπρστυφχψω",
			Length:   32,
		},
		nil,
	},
	{
		nil, `
url: example.com
length: "42"`,
		&pwclip.PWM{
			URL:     "example.com",
			Charset: pwclip.CharsetAlphaNumeric,
			Length:  42,
		},
		nil,
	},
	{
		nil, `
url: example.com
length: !!int 42`,
		&pwclip.PWM{
			URL:     "example.com",
			Charset: pwclip.CharsetAlphaNumeric,
			Length:  42,
		},
		nil,
	},
	{
		nil, `
url: "# Hello 世界"
username: \"\x41"\\\n
prefix: \☃foo"
charset: "\"\\\x41\u0042\U00000043 ☃\u2603\n"
length: 42`,
		&pwclip.PWM{
			URL:      "# Hello 世界",
			Username: `\"\x41"\\\n`,
			Prefix:   `\☃foo"`,
			Charset:  "\"\\ABC ☃\u2603\n",
			Length:   42,
		},
		nil,
	},
	{
		intptr(3), `
url: example
username: example@example.com
q1: frequent flier number
q2: first car`,
		nil,
		errors.New("question \"q3\" not in settings"),
	},
	{nil, "", nil, errors.New("unexpected yaml structure")},
	{nil, "[a,b,c]", nil, errors.New("unexpected yaml structure")},
	{nil, "url: [a,b,c]", nil, errors.New("url must be a string")},
	{nil, "length: 32x", nil, errors.New("length must be an int")},
}

func TestSettings(t *testing.T) {
	for i, test := range settingsTests {
		pwm, err := newPWMFromYaml([]byte(test.yaml), test.question)
		if !equalError(err, test.err) {
			t.Errorf("test %d (err):\n\texpected: %#v\n\tactually: %#v", i, test.err, err)
		}
		if !equalPWM(pwm, test.pwm) {
			t.Errorf("test %d (pwm):\n\texpected: %#v\n\tactually: %#v", i, test.pwm, pwm)
		}
	}
}

func equalError(e1 error, e2 error) bool {
	if e1 == e2 {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	return e1.Error() == e2.Error()
}

func equalPWM(p1 *pwclip.PWM, p2 *pwclip.PWM) bool {
	if p1 == p2 {
		return true
	}

	if p1 == nil || p2 == nil {
		return false
	}

	if p1.URL != p2.URL || p1.Username != p2.Username || p1.Prefix != p2.Prefix {
		return false
	}

	if p1.Length != p2.Length || !reflect.DeepEqual(p1.Charset, p2.Charset) {
		return false
	}

	e1 := p1.Extra
	e2 := p2.Extra
	if !((e1 == nil && e2 == nil) || (e1 != nil && e2 != nil && *e1 == *e2)) {
		return false
	}

	return true
}

func intptr(i int) *int       { return &i }
func strptr(s string) *string { return &s }
