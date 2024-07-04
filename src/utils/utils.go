package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	domain_errors "test-assignment/domain/errors"
	"test-assignment/domain/word"

	"github.com/fatih/color"
)

func RunGameLoop(secretWord string, count, maxLength int, errChan chan error) {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	// Print game  rules.
	PrintRules(count)

	// Initialize Reader once and re-use it in the loop.
	reader := bufio.NewReader(os.Stdin)

	var attempt int
	for attempt = 0; attempt < count; attempt++ {
		fmt.Printf("Guess the word (%d/%d): \n", attempt+1, maxLength)

		msg, err := reader.ReadString('\n')
		if err != nil && err != io.EOF { // Ignore the EOF error when read operation ends.
			fmt.Println(err)
			errChan <- err
			continue
		}

		// Remove any surrounding whitespaces including the newline.
		// Normalize the word by bringing it to uppercase.
		msg = strings.ToUpper(strings.TrimSpace(msg))

		// Verify the user's input content.
		err = word.IsValidUnicode(msg)
		if errors.Is(err, domain_errors.InvalidWordData{}) {
			attempt-- // If there was an input error, rollback the attempt count.
			fmt.Println(err)
			errChan <- err
			continue
		}

		// Verify the user's input length.
		err = word.IsValidLength(msg, maxLength)
		if errors.Is(err, domain_errors.InvalidWordLength{}) {
			attempt-- // If there was an input error, rollback the attempt count.
			fmt.Println(err)
			errChan <- err
			continue
		}

		if secretWord == msg {
			fmt.Println("You have won!!!")
			errChan <- nil
		} else {

			// Iterate through the secret word and user's input at same index to find GREEN matches.
			var letterIndx int
			for letterIndx = 0; letterIndx < len(secretWord); letterIndx++ {
				// e.g., Earth and Event both match on first index, the first E is one-to-one match and is GREEN colored.
				if secretWord[letterIndx] == msg[letterIndx] {
					green.Println("GREEN match" + " " + string(secretWord[letterIndx]) + " : " + string(msg[letterIndx]))
				}

				// Iterate through all leftover word and user's input letters to find YELLOW matches.
				var guessIndx int
				for guessIndx = 0; guessIndx < len(msg); guessIndx++ {
					if secretWord[letterIndx] == msg[letterIndx] {
						// It's the GREEN match scenario from above, skip this one.
						continue
					} else if secretWord[letterIndx] == msg[guessIndx] {
						// e.g., EarTh and EvenT both share E and T, the first E is one-to-one match and is GREEN match, T is YELLOW match.
						yellow.Println("YELLOW match" + " " + string(secretWord[letterIndx]) + " : " + string(msg[guessIndx]))
					}
				}
			}
		}
	}

	fmt.Println("Game end!")
	errChan <- nil
}

func PrintRules(count int) {
	fmt.Printf("Hello, this is a mini version of the web game Wordle in Go. \n\n")
	fmt.Printf("You have to guess the word and you have %v attempts. \n", count)
	fmt.Printf("Any letters that are in the right position" +
		"are highlighted in green while letters that are in the word but not in the" +
		"correct position will get a yellow outline. \n\n")
	fmt.Println("All words are 5 letter long, you can enter any word.")
}
