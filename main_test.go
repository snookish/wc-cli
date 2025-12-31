package main_test

import (
	"strings"
	"testing"

	counter "github.com/iamBelugax/wc-cli"
)

func TestCountWords(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{
			name:  "5 Words",
			input: "one two three four five",
			wants: 5,
		},
		{
			name:  "Empty Input",
			input: "",
			wants: 0,
		},
		{
			name:  "Single Space",
			input: " ",
			wants: 0,
		},
		{
			name:  "Single New Line",
			input: "one two three\nfour five",
			wants: 5,
		},
		{
			name:  "Multi New Lines",
			input: "one two three \n four \n five",
			wants: 5,
		},
		{
			name:  "Multi Spaces",
			input: "This is a sentence.  This is another one.",
			wants: 8,
		},
		{
			name:  "Prefixed Multi Spaces",
			input: "    This is a sentence.  This is another one.",
			wants: 8,
		},
		{
			name:  "Suffixed Multi Spaces",
			input: "This is a sentence.  This is another one.    ",
			wants: 8,
		},
		{
			name:  "Tab Character",
			input: "This is\ta sentence.\tThis is \tanother one.",
			wants: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			res := counter.CountWords(reader)
			if res != tc.wants {
				t.Logf("expected %d, got %d", tc.wants, res)
				t.Fail()
			}
		})
	}
}
