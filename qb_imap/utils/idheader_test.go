package utils_test

import (
	"testing"

	"github.com/rskvp/qb-lib/qb_imap/utils"
)

func TestFromIDHeader(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", ""},
		{"<>", ""},
		{"<%🤯>", "%🤯"},
		{"<foo@inbucket.org>", "foo@inbucket.org"},
		{"<foo%25bar>", "foo%bar"},
		{"foo+bar", "foo bar"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := utils.FromIDHeader(tc.input)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestToIDHeader(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", "<>"},
		{"foo@inbucket.org", "<foo@inbucket.org>"},
		{"foo%bar", "<foo%25bar>"},
		{"foo bar", "<foo+bar>"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := utils.ToIDHeader(tc.input)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
