package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
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

	sort.Slice(resp.Entities, func(i, j int) bool {
		return resp.Entities[i].Salience > resp.Entities[j].Salience
	})

	for _, entity := range resp.Entities {
		fmt.Printf("Name: %s, Salience: %f\n", entity.Name, entity.Salience)
	}

	return nil
}

func main() {
	url := validateURL("https://weareher.com/trans-dating/")
	parsedHTML := fetchContent(url)
	if err := analyzeEntities(parsedHTML); err != nil {
		log.Fatalf("Failed to analyse entities: %v", err)
	}

}
