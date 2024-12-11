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
			input:    "  dev     woRld Go  ",
			expected: []string{"dev", "world", "go"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected){
			t.Errorf("number of words does not match.\nExpected:\t%d\nActual\t%v\n", len(c.expected), len(actual))
			
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord{
				t.Errorf("word does not match:\nExpected:\t%v\nActual:\t%v\n", expectedWord, word)
			}
		}
	}
}

