package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	domain_errors "test-assignment/domain/errors"
	"test-assignment/domain/word"

	"github.com/fatih/color"
)

func RunGameLoop(count, maxLength int, errChan chan error) {
	green := color.New(color.FgGreen)

	// Print game  rules.
	PrintRules(count)

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

		msg = strings.TrimSpace(msg) // Remove any surrounding whitespace including the newline.

		err = word.IsValidUnicode(msg)
		if errors.Is(err, domain_errors.InvalidWordData{}) {
			fmt.Println(err)
			errChan <- err
			continue
		}

		err = word.IsValidLength(msg, maxLength)
		if errors.Is(err, domain_errors.InvalidWordLength{}) {
			fmt.Println(err)
			errChan <- err
			continue
		}

		green.Println("Word is fine, here we process it")
	}

	log.Println("Game end!")
	errChan <- nil
}

func PrintRules(count int) {
	fmt.Printf("Hello, this is a mini version of the web game Wordle in Go. \n\n")
	fmt.Printf("You have to guess the word and you have %v attempts. \n", count)
	fmt.Printf("Any letters that are in the right position" +
		"are highlighted in green while letters that are in the word but not in the" +
		"correct position will get a yellow outline \n\n")
	fmt.Println("All words are 5 letter long, you can enter any word.")
}
