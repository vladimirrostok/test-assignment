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
var color1 = color.New(color.FgGreen)
var color2 = color.New(color.FgYellow)

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
	// Print all the game rules.
	PrintRules(count)

	// Initialize Reader once and re-use it in the loop.
	reader := bufio.NewReader(os.Stdin)

	// Run the whole game in the loop until there are no more attempts or user wins the game.
	var attempt int
	for attempt = 0; attempt < count; attempt++ {
		fmt.Printf("\nGuess the word (%d/%d): \n", attempt+1, maxLength)

		// Read the user input.
		msg, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err // It's up to main process to decide whether game should continue when the when input is corrupted.
		}

		// Remove any surrounding whitespaces including the newline.
		// Normalize the word by bringing it to uppercase too.
		msg = strings.ToUpper(strings.TrimSpace(msg))

		// Verify the user's input data.
		err = word.IsValidUnicode(msg)
		if errors.Is(err, domain_errors.InvalidWordData{}) {
			attempt--        // If there was an input error, rollback the attempt count.
			fmt.Println(err) // Process game event in the real time.
			errChan <- err   // Send the game event to main process as an alternative (async) approach.
			continue         // Start new round.
		}

		// Verify the user's input length.
		err = word.IsValidLength(msg, maxLength)
		if errors.Is(err, domain_errors.InvalidWordLength{}) {
			attempt--        // If there was an input error, rollback the attempt count.
			fmt.Println(err) // Process game event in the real time.
			errChan <- err   // Send the game event to main process as an alternative (async) approach.
			continue         // Start new round.
		}

		// Check if guessed word matches the secret word.
		if secretWord == msg {
			color1.Println(msg)
			fmt.Println("You win!!!")
			errChan <- nil // Use existing channel to notify the main process.
			return         // Break the game loop from game itself and end there.
		} else {
			iterateWordMatches(secretWord, msg) // Run the complete guess-logic for this round.
		}
	}

	// End of the loop, if there are no attempts then show "Game end" and signal to end the game.
	fmt.Println("\nGame over!")
	errChan <- nil
}

// Game "guess the word" logic processor.
func iterateWordMatches(secretWord, msg string) string {
	// Hence we must know GREEN occurencies before we can accurately locate YELLOW instances,
	// we must process the entire string at least once, good example is water and otter = otTER.
	// We encode this in a few steps, 0 is "GREEN" and "1" is "YELLOW" instance.
	// We store entire calcultion progress in the variable and properly decode it in the end.
	var encodedResult = msg

	// Iterate through the secret word and user's input at same index to find GREEN matches.
	// We take full user input and iterate over the full secret word once finding [0]=[0] like matches.
	var i int
	for i = 0; i < len(msg); i++ {
		// e.g., Earth and Event both match on first index, the first E is one-to-one match and is GREEN colored.
		if msg[i] == secretWord[i] {
			encodedResult = replaceAtIndex(encodedResult, "0", i)
		}
	}

	// Iterate through the secret word and user's input to find YELLOW matches, we encode YELLOW into "1".
	for i := 0; i < len(msg); i++ {
		// Game rule says when guessed word has more instances of YELLOW letters than the secret word, then we don't make excessive yellow-s.
		// In case of water and otter we should highlight only otTER as green and ignore first one, not making the first "t" yellow.
		letterCountInSecretWord := strings.Count(secretWord, string(msg[i]))
		letterCountInUserInput := strings.Count(msg, string(msg[i]))

		// In case we have some instances left not covered with GREEN matches, looking for YELLOW instances.
		if letterCountInUserInput <= letterCountInSecretWord {
			var s int
			// Iterate every letter of word and input to find leftover matches.
			for s = 0; s < len(secretWord); s++ {
				// Dropping the already encoded values, 0 means GREEN match, GREEN matches are already pre-calculated and skipped.
				if encodedResult[i] == 0 {
					break
				} else if encodedResult[i] == secretWord[s] {
					// Encoding the new matches into the previous calculation result.
					encodedResult = replaceAtIndex(encodedResult, "1", i)
				}
			}
		}
	}

	printResult(encodedResult, msg)

	return encodedResult
}

func printResult(encodedResult, msg string) {
	// Visualize the result, decoding and printing the result.
	for i := 0; i < len(encodedResult); i++ {
		var encodedResult = string(encodedResult[i])
		var originalLetter = string(msg[i])

		// Paint values based on the encoding.
		if encodedResult == "0" {
			color1.Print(originalLetter)
		} else if encodedResult == "1" {
			color2.Print(originalLetter)
		} else {
			fmt.Print(originalLetter)
		}
	}
}

func replaceAtIndex(input, replacement string, index int) string {
	return input[:index] + string(replacement) + input[index+1:]
}
