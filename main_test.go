package main

import "testing"

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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := CountWords([]byte(tc.input))
			if res != tc.wants {
				t.Logf("expected %d, got %d", tc.wants, res)
				t.Fail()
			}
		})
	}
}
