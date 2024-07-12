package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"test-assignment/wordle/internal"
	"time"

	"github.com/fatih/color"
)

func main() {
	const wordGuesses = 5
	const wordLength = 5

	// Set flags to always print the line where the error came from.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load the words from file content.
	words, err := internal.ReadWordConfiguration("./config/words.txt")
	if err != nil {
		log.Fatalf("Fatal error %v \n", err)
	}

	// Create a new random number generator with a custom seed (e.g., current time).
	rndSource := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rndSource)

	// Pick any random word for the game session.
	randIndx := rng.Intn(len(words))               // We take the slice length and take any index from existing words.
	randomWord := strings.ToUpper(words[randIndx]) // Normalize the word which came from txt file, always turn it to uppercase.

	// Uncomment to verify the data.
	fmt.Println(words)
	fmt.Println("Randomly selected slice value : ", randomWord)

	// Provide channel for OS process termination signals.
	signalChan := make(chan os.Signal, 1)

	// Provide channel to catch errors from the game.
	errChan := make(chan error)

	// Listen to the OS termination signals.
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	// Always print a well-formatted message when the game ends.
	defer func() {
		var red = color.New(color.FgRed)
		var boldRed = red.Add(color.Bold)
		boldRed.Println("*** Thank you for playing the game! ***")
	}()

	// Run the game and listen to channels asynchronously unblocking the main thread.
	go internal.RunGameLoop(randomWord, wordGuesses, wordLength, errChan)

	// A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
	// Select statement without a "default" case is a blocking operation, so we don't really need extra loop around it here.
	select {
	case err := <-errChan:
		if err == nil {
			// When channel is closed by sender we receive nil and return from select back to main and end the app.
			// It is a "game end" event, so we use this instead of extra channel logic to catch this.
			// Technically, select exists at this point so we can leave this case empty, or use return/break here.
		} else if err == io.EOF {
			// When we execute the game in a Docker in non-interactive CLI we will get the EOF input error.
			log.Fatalf("Please don't run this in non-interactive CLI environment like Docker, error: %v", err)
		} else {
			// If there is some unexpected error, exit the game fast skipping any deferred operations, there is no recovery.
			log.Fatalf("Fatal error %v \n", err)
		}
	case <-signalChan: // If there is a shutdown signal, let game exit on its own.
		fmt.Print("Shutting down ...\n")
	}
}
