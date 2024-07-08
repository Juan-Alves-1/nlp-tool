package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
)

func fetchContent(url string) string {
	// Send an HTTP GET request
	allcontent, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch the URL: %v", url)
	}
	defer allcontent.Body.Close()

	// Read the response body
	bodycontent, err := io.ReadAll(allcontent.Body)
	if err != nil {
		fmt.Printf("failed to read the response body: %v", url)
	}

	// Return the HTML body as a string
	return string(bodycontent)

}

/*func validateUrl() (string, err) {
	// Check if a URL is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <URL>")
		return
	}

	// Get the URL from the command-line arguments
	url := os.Args[1]

	// Fetch the content
	approvedUrl, err := fetchContent(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Store the body content for later use (example: printing it here)
	return approvedUrl
}  */

// analyzeEntities sends a string of text to the Cloud Natural Language API to
// detect the entities of the text.
func analyzeEntities(html string) error {
	ctx := context.Background()

	// Initialize client.
	client, err := language.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	resp, err := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: html,
			},
			Type: languagepb.Document_HTML,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		return fmt.Errorf("AnalyzeEntities: %w", err)
	}
	fmt.Println(resp)
	return nil
}

func main() {
	parsedHTML := fetchContent("https://weareher.com/trans-dating")
	analyzeEntities(parsedHTML)

}
