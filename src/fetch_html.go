package source

import (
	"fmt"
	"io"
	"net/http"
)

func FetchContent(url string) (string, error) {

	allcontent, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch the URL: %v", url)
	}
	defer allcontent.Body.Close()

	bodyContent, err := io.ReadAll(allcontent.Body)
	if err != nil {
		fmt.Printf("failed to read the response body: %v", url)
		return "", err
	}

	// Most types of web pages, including Gutenberg content format - double-check content with: fmt.Printf("%s", bodyContent)
	return string(bodyContent), nil

	// For the traditional blog post format
	/* doc, err := readability.NewDocument(string(bodyContent))
	if err != nil {
		fmt.Println("Error parsing the HTML content:", err)
		return "", err
	}

	cleanHTML := doc.Content()
	return cleanHTML, nil */
}
