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
	"github.com/gocolly/colly"
)

type EntityInfo struct {
	Name        string
	Type        string
	Salience    float32
	HasWiki     bool
	WikiURL     string // via google API
	MentionedAs string
	MentionType languagepb.EntityMention_Type
}

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
	// fmt.Printf("%s", bodyContent)
	return string(bodyContent), nil

	// For traditional blog post format
	/* doc, err := readability.NewDocument(string(bodyContent))
	if err != nil {
		fmt.Println("Error parsing the HTML content:", err)
		return "", err
	}

	cleanHTML := doc.Content()
	return cleanHTML, nil */
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

	var entityInfos []EntityInfo

	for _, entity := range resp.Entities {
		hasWiki := false
		wikiURL := ""
		if val, ok := entity.Metadata["wikipedia_url"]; ok {
			hasWiki = true
			wikiURL = val
		}

		for _, mention := range entity.Mentions {
			entityInfos = append(entityInfos, EntityInfo{
				Name:        entity.Name,
				Type:        entity.Type.String(),
				Salience:    entity.Salience,
				HasWiki:     hasWiki,
				WikiURL:     wikiURL,
				MentionedAs: mention.Text.Content,
				MentionType: mention.Type,
			})
		}
	}

	// Sort entityInfos by salience in descending order
	sort.Slice(entityInfos, func(i, j int) bool {
		return entityInfos[i].Salience > entityInfos[j].Salience
	})

	// Print the 20 most prevalent entities
	fmt.Println("Top 20 Entities by Salience:")
	uniqueEntities := make(map[string]bool)
	count := 0
	for _, entity := range entityInfos {
		if count >= 20 {
			break
		}
		if !uniqueEntities[entity.Name] {
			uniqueEntities[entity.Name] = true
			fmt.Printf("Name: %s, Type: %s, Salience: %.6f, Has Wikipedia URL metadata: %t\n", entity.Name, entity.Type, entity.Salience, entity.HasWiki)
			count++
		}
	}

	// Print entities with mentions of type PROPER
	fmt.Println("\n_________________________________")
	fmt.Println("\nEntities with PROPER mentions:")
	var properEntities []EntityInfo
	for _, entity := range entityInfos {
		if entity.MentionType == languagepb.EntityMention_PROPER {
			properEntities = append(properEntities, entity)
		}
	}

	// Sort properEntities by salience in descending order
	sort.Slice(properEntities, func(i, j int) bool {
		return properEntities[i].Salience > properEntities[j].Salience
	})

	for _, entity := range properEntities {
		fmt.Printf("Name: %s, Text Content: %s, Salience: %6f, Wikipedia URL: %s\n", entity.Name, entity.MentionedAs, entity.Salience, entity.WikiURL)
	}

	return nil
}

func main() {
	url, err := validateURL("https://weareher.com/trans-dating")
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

	c := colly.NewCollector()

	var topResultURLs []string

	// Find and visit the first 10 links in the search results
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/url?q=") && len(topResultURLs) < 10 {
			url := strings.Split(href, "&")[0][7:] // Remove "/url?q=" and everything after "&"
			topResultURLs = append(topResultURLs, url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("___________________")
		fmt.Println("Visiting", r.URL)
	})

	err = c.Visit("https://www.google.com/search?q=trans+dating&hl=en&gl=us")
	if err != nil {
		log.Fatalf("Failed to visit Google search page: %v", err)
	}

	fmt.Println("Top 10 Google result URLs:")
	for i, url := range topResultURLs {
		fmt.Printf("%d: %s\n", i+1, url)
	}

}
