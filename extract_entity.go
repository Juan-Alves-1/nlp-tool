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
	"github.com/mauidude/go-readability"
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

func fetchContent(url string) (string, error) {

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

	// For Gutenberg content format
	/* fmt.Printf("%s", bodyContent)
	return string(bodyContent), nil */

	// For traditional blog post format
	doc, err := readability.NewDocument(string(bodyContent))
	if err != nil {
		fmt.Println("Error parsing the HTML content:", err)
		return "", err
	}

	cleanHTML := doc.Content()
	return cleanHTML, nil
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

	/* for _, mention := range entity.Mentions {
		fmt.Println(mention.Text)
	} */

	/* type myEntities struct {
		Name string
		Salience float32
		Wikipedia bool
		WikiUrl string
	}*/

	// Sort the entities by frequency in descending order
	sort.Slice(resp.Entities, func(i, j int) bool {
		return resp.Entities[i].Salience > resp.Entities[j].Salience
	})

	for _, entity := range resp.Entities[0:11] {
		fmt.Printf("Entity name: %s ", entity.Name)
		fmt.Printf("Entity type: %s ", entity.Type)
		fmt.Printf("Salience: %f ", entity.Salience)
		fmt.Println("______________________")
		for key, value := range entity.Metadata {
			fmt.Println(key)
			fmt.Println(value)
			fmt.Println("______________________")
			// fmt.Printf("Name: %s, Salience: %f\n", entity.Name, entity.Salience)
		}
	}
	return nil
}

func main() {
	url, err := validateURL("https://weareher.com/trans-dating/")
	if err != nil {
		log.Fatalf("Invalid URL: %s", err)
	}

	htmlContent, err := fetchContent(url)
	if err != nil {
		log.Fatalf("Failed to fetch: %s", err)
	}

	if err := analyzeEntities(htmlContent); err != nil {
		log.Fatalf("Failed to analyse entities: %v", err)
	}

}
