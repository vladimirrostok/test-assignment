package internal

import (
	"bufio"
	"errors"
	"os"
)

func ReadWordConfiguration(path string) ([]string, error) {
	var words []string

	if path == "" {
		return nil, errors.New("file route is missing")
	}

	file, err := os.Open(path)
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

	if len(words) == 0 {
		return nil, errors.New("please check the configuration words.txt file, it might be empty")
	}

	return words, nil
}
