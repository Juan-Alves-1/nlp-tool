package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
)

type Entity struct {
	Name            string
	Type            string
	Salience        float32
	WikiURLmetadata string // via Cloud Natural Language API
	WikiURLfromWiki string // via Wikipedia API
	MentionedAs     []string
}

func checkWikiURLfromWiki(entityName string) string {
	wikiURLfromWiki := "https://en.wikipedia.org/wiki/" + strings.ReplaceAll(entityName, " ", "_")

	resp, err := http.Get(wikiURLfromWiki)
	if err != nil {
		fmt.Printf("Failed to fetch %s", wikiURLfromWiki)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPermanentRedirect {
		return wikiURLfromWiki
	}
	return ""
}

func AnalyzeEntities(html string) ([]Entity, error) {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Couldn't create a context.Brackgroound: %s", err)
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
		log.Fatalf("Failed to analyse entities: %v", err)
	}

	var entityList []Entity

	for _, entity := range resp.Entities {
		e := Entity{
			Name:     entity.Name,
			Type:     entity.Type.String(),
			Salience: entity.Salience,
		}

		if url, ok := entity.Metadata["wikipedia_url"]; ok {
			e.WikiURLmetadata = url
		}

		for _, mention := range entity.Mentions {
			e.MentionedAs = append(e.MentionedAs, mention.String())
		}

		entityList = append(entityList, e)
	}

	// Sort entityInfos by salience in descending order
	sort.Slice(entityList, func(i, j int) bool {
		return entityList[i].Salience > entityList[j].Salience
	})

	// Delete duplicates
	uniqueEntities := make(map[string]bool)
	count := 0
	topEntities := []Entity{}
	for _, entity := range entityList {
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

	return topEntities, nil
}
