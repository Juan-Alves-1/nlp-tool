package source

import (
	"fmt"
)

func UrlInput() (string, error) {
	fmt.Println("What's the page you would like to analyse?")
	var urlInput string
	fmt.Scanln(&urlInput)
	return urlInput, nil
}
