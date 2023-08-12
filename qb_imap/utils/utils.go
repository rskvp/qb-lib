package utils

import (
	"bytes"
	"net/mail"
	"unicode"

	qbl "github.com/rskvp/qb-lib"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var latinSpecialMap = map[rune]rune{
	'Đ': 'D',
	'đ': 'd',
	'Ħ': 'H',
	'ħ': 'h',
	'ĸ': 'K',
	'Ŀ': 'L',
	'Ł': 'L',
	'ŉ': 'n',
	'Ŋ': 'N',
	'ŋ': 'n',
	'Ŧ': 'T',
	'ŧ': 't',
}

func Html2TextFromString(input string) (string, error) {
	parser, err := qbl.HTML.NewParser(input)
	if nil != err {
		return "", err
	}
	return parser.TextAll(), nil
}

// Wrap builds a byte slice from strs, wrapping on word boundaries before max chars
func Wrap(max int, strs ...string) []byte {
	input := make([]byte, 0)
	output := make([]byte, 0)
	for _, s := range strs {
		input = append(input, []byte(s)...)
	}
	if len(input) < max {
		// Doesn't need to be wrapped
		return input
	}
	ls := -1 // Last seen space index
	lw := -1 // Last written byte index
	ll := 0  // Length of current line
	for i := 0; i < len(input); i++ {
		ll++
		switch input[i] {
		case ' ', '\t':
			ls = i
		}
		if ll >= max {
			if ls >= 0 {
				output = append(output, input[lw+1:ls]...)
				output = append(output, '\r', '\n', ' ')
				lw = ls // Jump over the space we broke on
				ll = 1  // Count leading space above
				// Rewind
				i = lw + 1
				ls = -1
			}
		}
	}
	return append(output, input[lw+1:]...)
}

// JoinAddress formats a slice of Address structs such that they can be used in a To or Cc header.
func JoinAddress(address []mail.Address) string {
	if len(address) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	for i, a := range address {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(a.String())
	}
	return buf.String()
}

// ToASCII converts unicode to ASCII by stripping accents and converting some special characters
// into their ASCII approximations.  Anything else will be replaced with an underscore.
func ToASCII(s string) string {
	// unicode.Mn: nonspacing marks
	tr := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), runes.Map(mapLatinSpecial),
		norm.NFC)
	r, _, _ := transform.String(tr, s)
	return r
}

// mapLatinSpecial attempts to map non-accented latin extended runes to ASCII
func mapLatinSpecial(r rune) rune {
	if v, ok := latinSpecialMap[r]; ok {
		return v
	}
	if r > 0x7e {
		return '_'
	}
	return r
}
