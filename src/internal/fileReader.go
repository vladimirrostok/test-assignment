package internal

import (
	"bufio"
	"os"
)

func ReadWordConfiguration() ([]string, error) {
	var words []string

	file, err := os.Open("./configuration/words.txt")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)

	for Scanner.Scan() {
		words = append(words, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
