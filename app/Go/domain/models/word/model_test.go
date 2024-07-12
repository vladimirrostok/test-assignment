package word

import (
	"errors"
	domain_errors "test-assignment/domain/errors"
	"testing"
)

type wordDataValidationTest struct {
	word     string
	expected error
}

type wordLengthValidationTest struct {
	word     string
	length   int
	expected error
}

var wordDataValidationTests = []wordDataValidationTest{
	{"Word1", domain_errors.InvalidWordData{}},
	{"Wo rd", domain_errors.InvalidWordData{}},
	{"Word$", domain_errors.InvalidWordData{}},
	{"Words", nil},
}

var wordLengthValidationTests = []wordLengthValidationTest{
	{"W", 5, domain_errors.InvalidWordLength{}},
	{"Wor", 5, domain_errors.InvalidWordLength{}},
	{"Wordss", 5, domain_errors.InvalidWordLength{}},
	{"Words", 5, nil},
}

func TestWordDataValidation(t *testing.T) {
	for _, test := range wordDataValidationTests {
		output := IsValidUnicode(test.word)
		if !errors.Is(output, test.expected) {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}

func TestWordLengthValidation(t *testing.T) {
	for _, test := range wordLengthValidationTests {
		output := IsValidLength(test.word, test.length)
		if !errors.Is(output, test.expected) {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
