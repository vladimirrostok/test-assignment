package domain_errors

type (
	// InvalidData signifies a user input with forbidden non-unicode data.
	InvalidWordData struct{}

	// InvalidData signifies a user input of incorrect length.
	InvalidWordLength struct{}
)

func (err InvalidWordData) Error() string {
	return "Invalid data, non-unicode characters present"
}

func (err InvalidWordLength) Error() string {
	return "Invalid data, incorrect length"
}
