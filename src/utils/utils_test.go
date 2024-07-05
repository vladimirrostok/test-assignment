package utils

import (
	"fmt"
	"testing"
)

// NB!
// Keep in mind that result encoding is the root functionality of this game logic,
// thus we verify the resulted value is what we expect from the game.
// GREEN = 0, YELLOW =1, to simplify I put it in comments like GREEN and <yellow letter>.
type wordIterationTest struct {
	secretWord, msg string
	expected        string
}

var wordDataValidationTests = []wordIterationTest{
	{"WATER", "OTTER", "OT000"}, // it's otTER
	{"WATER", "TOTER", "TO000"}, // it's toTER
	{"WATER", "TTTTT", "TT0TT"}, // it's ttTtt
	{"HOUND", "HONDA", "0011A"}, // it's HO<ND>a
	{"HOUND", "DNUOH", "11011"}, // it's <HO>U<ND>
	{"HOUND", "DNOUH", "11111"}, // it's <HOUND>
	{"AAAAA", "BBBBB", "BBBBB"}, // it's bbbbb
	{"WATER", "WATER", "00000"}, // it's WATER
}

func TestWordIteration(t *testing.T) {
	for _, test := range wordDataValidationTests {
		if output := iterateWordMatches(test.secretWord, test.msg); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}

func TestWordDecodedPrint(t *testing.T) {
	// Quick printout to the console to see the colorful output.
	// Not automated solution, good utility to verify hardcoded cases manually.
	printResult("OT000", "OTTER")
	fmt.Println("")
	printResult("0011A", "HONDA")
	fmt.Println("")
}
