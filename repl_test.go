package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   Mixed  CASE    and   spacing  ",
			expected: []string{"mixed", "case", "and", "spacing"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("For input %q, expected %d words but got %d: %v",
				c.input, len(c.expected), len(actual), actual)
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("For input %q, word at index %d expected %q but got %q",
					c.input, i, c.expected[i], actual[i])
			}
		}
	}
}
