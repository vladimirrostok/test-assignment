package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	domain_errors "test-assignment/domain/errors"
	"test-assignment/domain/word"

	"github.com/fatih/color"
)

// Initialize the colors and re-use them in the loop.
var greenBackground = color.New(color.BgGreen)
var yellowBackground = color.New(color.BgYellow)

// We use this variable at two places, move outside of func to simplify the tests input.
var wordMaxLength int

// Print all the game rules.
func PrintRules(count int) {
	fmt.Printf("Hello, this is a mini version of the web game Wordle in Go. \n\n")
	fmt.Printf("You have to guess the word and you have %v attempts. \n", count)
	fmt.Printf("Any letters that are in the right position" +
		"are highlighted in green while letters that are in the word but not in the" +
		"correct position will get a yellow outline. \n\n")
	fmt.Println("All words are 5 letter long, you can enter any word.")
}

func RunGameLoop(secretWord string, count, maxLength int, errChan chan error) {
	defer func() {
		close(errChan) // As we know any processing has been finished, we can safely close chan from the sender.
	}()

	wordMaxLength = maxLength

	// Print all the game rules.
	PrintRules(count)

	// Initialize Reader once and re-use it in the loop.
	reader := bufio.NewReader(os.Stdin)

	// Run the whole game in the loop until there are no more attempts or user wins the game.
	var attempt int
	for attempt = 0; attempt < count; attempt++ {
		fmt.Printf("\nGuess the word (%d/%d): ", attempt+1, maxLength)

		// Read the user input.
		input, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err // It's up to main process to decide whether game should continue when the when input is corrupted.
		}

		// Remove any surrounding whitespaces including the newline.
		// Normalize the word by bringing it to uppercase too.
		input = strings.ToUpper(strings.TrimSpace(input))

		// Verify the user's input data.
		err = word.IsValidUnicode(input)
		if errors.Is(err, domain_errors.InvalidWordData{}) {
			attempt--        // If there was an input error, rollback the attempt count.
			fmt.Println(err) // Process game event in the real time.
			continue         // Start new round.
		}

		// Verify the user's input length.
		err = word.IsValidLength(input, maxLength)
		if errors.Is(err, domain_errors.InvalidWordLength{}) {
			attempt--        // If there was an input error, rollback the attempt count.
			fmt.Println(err) // Process game event in the real time.
			continue         // Start new round.
		}

		// Check if guessed word matches the secret word.
		if secretWord == input {
			greenBackground.Println(input)
			fmt.Println("You win!!!")
			return // Break the game loop from game itself and end there.
		} else {
			iterateWordMatches([]rune(secretWord), []rune(input)) // Run the complete guess-logic for this round.
		}
	}

	// End of the loop, if there are no attempts then show "Game end".
	fmt.Println("\nGame over!")
}

// Game "guess the word" logic processor.
func iterateWordMatches(secretWord, input []rune) string {
	// Hence we must know GREEN occurencies before we can accurately locate <yellow> instances,
	// we must process the entire word once, good tricky example case for that is: water and otter = o<t>TER.
	// We encode results as: 0 is "GREEN" and "1" is <yellow> instance, and decode it while printing it.

	// Hence we create new slice based on existing slice, we don't want to modify original slice, thus we create new slice.
	var encodedInput = make([]rune, len(input))
	copy(encodedInput, input)

	// For complex scenarios we need to track how much we used <yellow> letters before painting new ones.
	// Rule says - when guessed word has more instances of <yellow> letter than the secret word, then we don't make excessive <yellow>-s.
	// Example: "water" and "otter" should highlight only otTER, TER as GREEN and first t is ignored, not making the first "t" <yellow>.
	// Another example: otter and toott = <to>o<t>t (2 <yellow> T and 1 <yellow> O, no GREEN here).
	// First we calculate the easy GREEN matches, take this number and calculate the leftover <yellow> amount.
	type letterUsage struct {
		matchesTotal int
		matchesUsed  int
	}

	// We create map of fixed size as we already know the max possible amount of unique letters in the input.
	yellowsUsed := make(map[rune]*letterUsage, wordMaxLength)

	// Iterate through the secret word and user's input at same index to find GREEN matches.
	// We take full user input and iterate over the full secret word once finding [0]=[0] like GREEN matches.
	var i int
	for i = 0; i < len(input); i++ {
		letter := input[i]
		totalLettersInSecret := strings.Count(string(secretWord), string(letter)) // Take total amount of letter instances in word.

		// We fill the map with the secret word's letter as key and amount as value to use it after.
		// We need to do it only once, so we check if we added something here before, if exists just skip.
		if yellowsUsed[letter] == nil {
			yellowsUsed[letter] = &letterUsage{matchesTotal: totalLettersInSecret, matchesUsed: 0}
		}

		// e.g., Earth and Event both match on first index, the first E is one-to-one match and is GREEN colored.
		// We are encoding the color 0|1 right into the string, it's easy to iterate over it to paint proper positions,
		// so we just use this to carry the correct colored indexes and colors of these matched letters.
		if letter == secretWord[i] {
			encodedInput[i] = '0'
			yellowsUsed[letter].matchesUsed++
		}
	}

	// Iterate through the secret word and user's input to find <yellow> matches, we encode <yellow> into "1".
	// Game rule says when guessed word has more instances of <yellow> letters than the secret word, then we don't make excessive <yellow>.
	// In case of water and otter we should highlight only otTER as green and ignore first one, not making the first "t" yellow.
	// GREEN matches are already encoded and we just skip these cycles.
	for i := 0; i < len(input); i++ {
		letter := input[i]

		// In case we have some <yellow> matches left we look for it, and checking the amount too.
		// Iterate every letter of word and input to find leftover <yellow> matches, skipping [0]=[0] GREEN from previous step.
		var s int
		for s = 0; s < len(secretWord); s++ {
			// Dropping the already encoded values, 0 means GREEN match, GREEN matches are already pre-calculated and skipped.
			if encodedInput[i] == '0' {
				break // Go to next letter, this GREEN is simply skipped.
			} else if encodedInput[i] == secretWord[s] {
				// If amount of used yellow matches is less than total amount of possible matches then use this position (from left to right).
				if yellowsUsed[letter].matchesUsed < yellowsUsed[letter].matchesTotal {
					yellowsUsed[letter].matchesUsed++
					encodedInput[i] = '1'
				}
			}
		}
	}

	printResult(encodedInput, input) // Print the result for the entered word.

	return string(encodedInput) // Return the results as string simplyfing unit tests for this.
}

func printResult(encodedInput, input []rune) {
	// Visualize the result, decoding and printing the result.
	for i := 0; i < len(encodedInput); i++ {

		// Paint values based on the encoding.
		if encodedInput[i] == '0' {
			greenBackground.Print(string(input[i]))
		} else if encodedInput[i] == '1' {
			yellowBackground.Print(string(input[i]))
		} else {
			fmt.Print(string(input[i]))
		}
	}
}
