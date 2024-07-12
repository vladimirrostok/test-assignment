package main

import (
	"fmt"
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
	// Select statement without a "default" case is a blocking operation, so we don't really need "for" loop here.
	select {
	case err := <-errChan: // If there is unknown error, then let it exit fast.
		if err != nil { // When we close channel from the sender, we receive nil here.
			log.Fatalf("Fatal error %v \n", err) // Exit the game fast, skip processing any defer-calls.
		}
	case <-signalChan: // If there is a shutdown signal, let game exit on its own.
		fmt.Print("Shutting down ...\n")
		return // Exit the loop, there's nothing blocking after it so the app closes.
	}
}
