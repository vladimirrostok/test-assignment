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

func RunGameLoop(secretWord string, count, maxLength int, errChan chan error) {
	// Initialize the colors and re-use them in the loop.
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	// Print all the game rules.
	PrintRules(count)

	// Initialize Reader once and re-use it in the loop.
	reader := bufio.NewReader(os.Stdin)

	// Run the whole game in the loop until there are no more attempts or user wins the game.
	var attempt int
	for attempt = 0; attempt < count; attempt++ {
		fmt.Printf("Guess the word (%d/%d): \n", attempt+1, maxLength)

		// Read the user input.
		msg, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err // It's up to main process if game should continue when the when input is corrupted.
		}

		// Remove any surrounding whitespaces including the newline.
		// Normalize the word by bringing it to uppercase too.
		msg = strings.ToUpper(strings.TrimSpace(msg))

		// Verify the user's input data.
		err = word.IsValidUnicode(msg)
		if errors.Is(err, domain_errors.InvalidWordData{}) {
			attempt--      // If there was an input error, rollback the attempt count.
			errChan <- err // Send the game event to main process as an alternative (async) approach.
			continue       // Start new round.
		}

		// Verify the user's input length.
		err = word.IsValidLength(msg, maxLength)
		if errors.Is(err, domain_errors.InvalidWordLength{}) {
			attempt--        // If there was an input error, rollback the attempt count.
			fmt.Println(err) // Process game event in the real time.
			errChan <- err   // Send the game event to main process as an alternative (async) approach.
			continue         // Start new round.
		}

		// Check if guessed word matches the secret word fully.
		if secretWord == msg {
			fmt.Println("You have won!!!")
			errChan <- nil // Use existing channel to command the main process to end the game.
		} else {
			iterateWordMatches(secretWord, msg, green, yellow) // Run the complete guess-logic for this round.
		}
	}

	// End of the loop, if there are no attempts then show "Game end" and signal to end the game.
	fmt.Println("Game end!")
	errChan <- nil
}

// Print all the game rules.
func PrintRules(count int) {
	fmt.Printf("Hello, this is a mini version of the web game Wordle in Go. \n\n")
	fmt.Printf("You have to guess the word and you have %v attempts. \n", count)
	fmt.Printf("Any letters that are in the right position" +
		"are highlighted in green while letters that are in the word but not in the" +
		"correct position will get a yellow outline. \n\n")
	fmt.Println("All words are 5 letter long, you can enter any word.")
}

// Game "guess the word" logic processor.
func iterateWordMatches(secretWord, msg string, color1, color2 *color.Color) {

	// Iterate through the secret word and user's input at same index to find GREEN matches.
	var letterIndx int
	for letterIndx = 0; letterIndx < len(secretWord); letterIndx++ {
		// e.g., Earth and Event both match on first index, the first E is one-to-one match and is GREEN colored.
		if secretWord[letterIndx] == msg[letterIndx] {
			color1.Println("GREEN match" + " " + string(secretWord[letterIndx]) + " : " + string(msg[letterIndx]))
		}

		// Iterate through all leftover word and user's input letters to find YELLOW matches.
		var guessIndx int
		for guessIndx = 0; guessIndx < len(msg); guessIndx++ {
			if secretWord[letterIndx] == msg[letterIndx] {
				// It's the GREEN match scenario from above, skip this one.
				continue
			} else if secretWord[letterIndx] == msg[guessIndx] {
				// e.g., EarTh and EvenT both share E and T, the first E is one-to-one match and is GREEN match, T is YELLOW match.
				color2.Println("YELLOW match" + " " + string(secretWord[letterIndx]) + " : " + string(msg[guessIndx]))
			}
		}
	}
}
