package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"test-assignment/utils"
	"time"

	"github.com/fatih/color"
)

func main() {
	const wordGuesses = 5
	const wordLength = 5

	// Load the words file content.
	words, err := utils.ReadWordConfiguration()
	if err != nil {
		fmt.Printf("Fatal error %v \n", err)
	} else if len(words) == 0 {
		fmt.Println("Please check the configuration words.txt file, it might be empty")
		os.Exit(0)
	}

	// Create a new random number generator with a custom seed (e.g., current time).
	rndSource := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rndSource)

	// Pick any random word for the game session.
	randIndx := rng.Intn(len(words))
	var randomWord = words[randIndx]
	randomWord = strings.ToUpper(randomWord) // Normalize the word which came from txt file, always turn it to uppercase.

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

	// Process the game and final print.
	go func() {
		defer func() {
			// Always print a well-formatted message when the game ends.
			var red = color.New(color.FgRed)
			var boldRed = red.Add(color.Bold)
			boldRed.Println("*** Thank you for playing the game! ***")

			os.Exit(0)
		}()

		// Start the game.
		utils.RunGameLoop(randomWord, wordGuesses, wordLength, errChan)
	}()

	// Block until the error or forced shutdown signal kicks in.
	for {
		select {
		case err := <-errChan: // If there is unknown error, then let it exit fast.
			log.Printf("Fatal error %v \n", err)
			os.Exit(0) // Exit the game fast, skip processing any defer-calls.
		case <-signalChan: // If there is a shutdown signal, let game exit on its own.
			fmt.Print("Shutting down ...\n")
			return // Exit the loop, there's nothing blocking after it so the app closes.
		}
	}
}
