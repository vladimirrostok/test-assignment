package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	domain_errors "test-assignment/domain/errors"
	"test-assignment/utils"

	"github.com/fatih/color"
)

// Example words are: water, otter, hound, pizza, eagle, fruit, paper.
const wordGuesses = 5
const wordLength = 5

// Available words to guess.
//var words = [7]string{"water", "otter", "hound", "pizza", "eagle", "fruit", "paper"}

func main() {
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
		utils.RunGameLoop(wordGuesses, wordLength, errorsChan)
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
