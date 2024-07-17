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
	Name            string
	Type            string
	Salience        float32
	WikiURLmetadata string // via Cloud Natural Language API
	WikiURLfromWiki string // via Wikipedia API
	MentionedAs     string
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

func checkWikiURLfromWiki(entityName string) string {
	WikiURLfromWiki := "https://en.wikipedia.org/wiki/" + strings.ReplaceAll(entityName, " ", "_")

	resp, err := http.Get(WikiURLfromWiki)
	if err != nil {
		fmt.Printf("Failed to fetch %s", WikiURLfromWiki)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPermanentRedirect {
		return WikiURLfromWiki
	}
	return ""
}

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
		WikiURLmetadata := ""
		if key, ok := entity.Metadata["wikipedia_url"]; ok {
			WikiURLmetadata = key
		}

		for _, mention := range entity.Mentions {
			entityInfos = append(entityInfos, EntityInfo{
				Name:            entity.Name,
				Type:            entity.Type.String(),
				Salience:        entity.Salience,
				WikiURLmetadata: WikiURLmetadata,
				WikiURLfromWiki: "",
				MentionedAs:     mention.Text.Content,
			})
		}
	}

	// Sort entityInfos by salience in descending order
	sort.Slice(entityInfos, func(i, j int) bool {
		return entityInfos[i].Salience > entityInfos[j].Salience
	})

	// Delete duplicates
	uniqueEntities := make(map[string]bool)
	count := 0
	topEntities := []EntityInfo{}
	for _, entity := range entityInfos {
		if count >= 30 {
			break
		}
		if !uniqueEntities[entity.Name] {
			uniqueEntities[entity.Name] = true
			topEntities = append(topEntities, entity)
			count++
		}
	}

	// Check Wikipedia URLs for the top 30 entities
	for i := 0; i < len(topEntities); i++ {
		if topEntities[i].WikiURLmetadata == "" {
			topEntities[i].WikiURLfromWiki = checkWikiURLfromWiki(topEntities[i].Name)
		}
	}

	// Print the top 30 entities with Wikipedia URL information
	for _, entity := range topEntities {
		fmt.Printf("Name: %s, Type: %s, Salience: %.6f, Has Wikipedia URL metadata: %t, Has Wikipedia URL: %s\n",
			entity.Name, entity.Type, entity.Salience, entity.WikiURLmetadata != "", entity.WikiURLfromWiki)
	}

	return nil
}

func main() {
	url, err := validateURL("https://weareher.com/")
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

	// Find and visit the first 12 links in the search results
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/url?q=") && len(topResultURLs) < 12 {
			url := strings.Split(href, "&")[0][7:] // Remove "/url?q=" and everything after "&"
			topResultURLs = append(topResultURLs, url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("___________________")
		fmt.Println("Visiting", r.URL)
	})

	err = c.Visit("https://www.google.com/search?q=lesbian+dating&hl=en&gl=us")
	if err != nil {
		log.Fatalf("Failed to visit Google search page: %v", err)
	}

	fmt.Println("Top 10 Google result URLs:")
	for i, url := range topResultURLs[2:] {
		fmt.Printf("%d: %s\n", i+1, url)
	}

}
