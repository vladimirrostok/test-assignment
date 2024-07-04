package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	domain_errors "test-assignment/domain/errors"
	"test-assignment/utils"
	"time"

	"github.com/fatih/color"
)

const wordGuesses = 5
const wordLength = 5

func main() {
	// Load the words file content
	words, err := utils.ReadWordConfiguration()
	if err != nil {
		fmt.Printf("Fatal error %v \n", err)
	}

	// Create a new random number generator with a custom seed (e.g., current time)
	rndSource := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rndSource)
	randIndx := rng.Intn(len(words))
	var randomWord = words[randIndx]
	randomWord = strings.ToUpper(randomWord) // Normalize the word which came from txt file, always turn it to uppercase.
	fmt.Println(words)
	fmt.Println("Randomly selected slice value : ", randomWord)

	// Always print a well-formatted message when the game ends.
	// Combine defer call in a single function to print before Exit.
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	defer func() {
		boldRed.Println("*** Thank you for playing the game! ***")
		os.Exit(0)
	}()

	// Provide channel for errors from the process goroutines.
	errorsChan := make(chan error)

	go func() {
		// Start the game loop, re-use the initialized reader.
		utils.RunGameLoop(randomWord, wordGuesses, wordLength, errorsChan)
	}()

	// Provide channel for OS process termination signals.
	signalChan := make(chan os.Signal, 1)

	// Listen to the OS termination signals.
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	// Block till err/termination chan comes in.
	for {
		select {
		case err := <-errorsChan:
			if errors.Is(err, domain_errors.InvalidWordData{}) {
				// Here we can catch in-game event asynchronously.
				// log.Println(err)
			} else if errors.Is(err, domain_errors.InvalidWordLength{}) {
				// Here we can catch in-game event asynchronously.
				// log.Println(err)
			} else if err == nil {
				// Use this as a signal to quit the game.
				// TODO: Replace it with a proper channel.
				runtime.Goexit()
			} else {
				// If there is unknown error, then fail hard and let it exit.
				fmt.Printf("Fatal error %v \n", err)
				runtime.Goexit()
			}
		case <-signalChan:
			fmt.Print("Shutting down in a second...\n")
			runtime.Goexit()
		}
	}
}
