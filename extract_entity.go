package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
)

func validateURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "Invalid URL format"
	}
	validURL := parsedURL.String()
	return validURL
}

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
	url := validateURL("https://weareher.com/trans-dating/")
	parsedHTML := fetchContent(url)
	analyzeEntities(parsedHTML)

}
