package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"test-assignment/utils"
	"time"

	"github.com/fatih/color"
)

const wordGuesses = 5
const wordLength = 5

func main() {
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
	)

	// Listen to the game when it ends.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the game asynchronously.
	go utils.RunGameLoop(randomWord, wordGuesses, wordLength, errChan, &wg)

	// Process the game end.
	go func() {
		defer func() {
			// Always print a well-formatted message when the game ends.
			var red = color.New(color.FgRed)
			var boldRed = red.Add(color.Bold)
			boldRed.Println("*** Thank you for playing the game! ***")

			os.Exit(0)
		}()
		// Wait for all workers to complete.
		wg.Wait()
	}()

	// Block until the error or forced shutdown signal kicks in.
	for {
		select {
		// If there is unknown error, then let it exit fast.
		case err := <-errChan:
			log.Printf("Fatal error %v \n", err)
			os.Exit(0)
			// If there is a shutdown signal, let game exit on its own.
		case <-signalChan:
			fmt.Print("Shutting down ...\n")
			os.Exit(0)
		}
	}
}
