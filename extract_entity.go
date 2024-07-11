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

	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
	"github.com/PuerkitoBio/goquery"
)

func validateURL(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "Invalid URL format", err
	}
	validURL := parsedURL.String()
	return validURL, nil
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

func extractBodyContent(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	bodyHtml, err := doc.Find("body").Html()
	if err != nil {
		return "", fmt.Errorf("failed to extract body content: %v", err)
	}

	return bodyHtml, nil
}

// not working
/* func extractContentFromFirstH1(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract the content of the <body> tag
	bodySelection := doc.Find("body")
	if bodySelection.Length() == 0 {
		return "", fmt.Errorf("no <body> tag found")
	}
	bodyContent, err := bodySelection.Html()
	if err != nil {
		return "", fmt.Errorf("failed to get <body> content: %v", err)
	}

	// Parse the body content again to find the first <h1> tag and extract content from there onward
	bodyDoc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse body content: %v", err)
	}

	var contentBuilder strings.Builder
	startAppending := false

	bodyDoc.Find("body").Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("h1") {
			startAppending = true
		}
		if startAppending {
			html, err := s.Html()
			if err == nil {
				contentBuilder.WriteString(html)
			} else {
				fmt.Printf("Error parsing HTML: %v\n", err)
			}
		}
	})

	content := contentBuilder.String()
	if content == "" {
		return "", fmt.Errorf("no content found after the first <h1> tag")
	}
	return content, nil
} */

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

	// Create a map to store the frequency of each entity
	entityFrequency := make(map[string]int)

	for _, entity := range resp.Entities {
		entityFrequency[entity.Name] += len(entity.Mentions)
	}

	// Create a slice to sort entities by frequency
	type entityInfo struct {
		Name      string
		Type      string
		Frequency int
	}
	var entities []entityInfo

	for _, entity := range resp.Entities {
		entities = append(entities, entityInfo{
			Name:      entity.Name,
			Type:      entity.Type.String(),
			Frequency: entityFrequency[entity.Name],
		})
	}

	// Sort the entities by frequency in descending order
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Frequency > entities[j].Frequency
	})

	// Print the sorted entities
	for _, entity := range entities {
		fmt.Printf("Name: %s, Type: %s, Frequency: %d\n", entity.Name, entity.Type, entity.Frequency)
	}

	return nil
}

func main() {
	url, err := validateURL("https://weareher.com/trans-dating/")
	if err != nil {
		log.Fatalf("Invalid URL: %s", err)
	}
	htmlContent := fetchContent(url)

	bodyContent, err := extractBodyContent(htmlContent)
	if err != nil {
		log.Fatalf("Failed to extract body content: %v", err)
	}

	if err := analyzeEntities(bodyContent); err != nil {
		log.Fatalf("Failed to analyse entities: %v", err)
	}

}
