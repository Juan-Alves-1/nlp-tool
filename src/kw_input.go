package source

import (
	"bufio"
	"fmt"
	"os"
)

func KwInput() (string, error) {
	fmt.Println("What's your target keyword?")

	// Use bufio.NewReader to read the entire line of input
	reader := bufio.NewReader(os.Stdin)
	kwInput, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}

	return kwInput, nil
}
