package utils_test

import (
	"fmt"
	"net/mail"
	"testing"
	"unicode/utf8"

	"github.com/rskvp/qb-lib/qb_imap/utils"
)

func TestJoinAddressEmpty(t *testing.T) {
	got := utils.JoinAddress(make([]mail.Address, 0))
	if got != "" {
		t.Errorf("Empty list got: %q, wanted empty string", got)
	}
}

func TestJoinAddressSingle(t *testing.T) {
	input := []mail.Address{
		{Name: "", Address: "one@bar.com"},
	}
	want := "<one@bar.com>"
	got := utils.JoinAddress(input)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	input = []mail.Address{
		{Name: "one name", Address: "one@bar.com"},
	}
	want = `"one name" <one@bar.com>`
	got = utils.JoinAddress(input)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestJoinAddressMany(t *testing.T) {
	input := []mail.Address{
		{Name: "one", Address: "one@bar.com"},
		{Name: "", Address: "two@foo.com"},
		{Name: "three", Address: "three@baz.com"},
	}
	want := `"one" <one@bar.com>, <two@foo.com>, "three" <three@baz.com>`
	got := utils.JoinAddress(input)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestToASCII(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", ""},
		{"Yoùr Śtring", "Your String"},
		{"šđčćžŠĐČĆŽŁĲſ", "sdcczSDCCZL__"},
		{"Ötzi's Nationalität èàì", "Otzi's Nationalitat eai"},
	}
	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			got := utils.ToASCII(tc.input)
			if got != tc.want {
				t.Errorf("Got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestToASCIILatinExt(t *testing.T) {
	b := make([]byte, 3)
	for r := rune(0x100); r <= 0x17f; r++ {
		n := utf8.EncodeRune(b, r)
		in := string(b[:n])
		out := utils.ToASCII(in)
		fmt.Printf("%c %q %q\n", r, in, out)
		if out == "" {
			t.Errorf("ToASCII(%q) returned empty string", in)
		}
		got, _ := utf8.DecodeRuneInString(out)
		if got == utf8.RuneError {
			t.Errorf("ToASCII(%q) returned undecodable rune: %q", in, out)
		}
		if got < 0x21 || 0x7e < got {
			t.Errorf("ToASCII(%q) returned non-ASCII rune: %c (%U)", in, got, got)
		}
	}
}

func TestWrapEmpty(t *testing.T) {
	b := utils.Wrap(80, "")
	got := string(b)

	if got != "" {
		t.Errorf(`got: %q, want: ""`, got)
	}
}

func TestWrapIdentityShort(t *testing.T) {
	want := "short string"
	b := utils.Wrap(15, want)
	got := string(b)

	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestWrapIdentityLong(t *testing.T) {
	want := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	b := utils.Wrap(5, want)
	got := string(b)

	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestWrap(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{
			"one two three",
			"one\r\n two\r\n three",
		},
		{
			"a bb ccc dddd eeeee ffffff",
			"a bb\r\n ccc\r\n dddd\r\n eeeee\r\n ffffff",
		},
		{
			"aaaaaa bbbbb cccc ddd ee f",
			"aaaaaa\r\n bbbbb\r\n cccc\r\n ddd\r\n ee f",
		},
		{
			"1 3 5 1 3 5 1 3 5",
			"1 3 5\r\n 1 3 5\r\n 1 3 5",
		},
		{
			"55555 55555 55555",
			"55555\r\n 55555\r\n 55555",
		},
		{
			"666666 666666 666666",
			"666666\r\n 666666\r\n 666666",
		},
		{
			"7777777 7777777 7777777",
			"7777777\r\n 7777777\r\n 7777777",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b := utils.Wrap(6, tc.input)
			got := string(b)
			if got != tc.want {
				t.Errorf("got: %q, want: %q", got, tc.want)
			}

		})
	}
}
