package main

import (
	"fmt"
	"log"
	source "nlp-tool/src"
)

func main() {
	/*// Get top entities for a given URL
	url, err := source.UrlInput()
	if err != nil {
		log.Fatalf("Invalid input: %s", err)
	}

	validURL, err := source.ValidateURL(url)
	if err != nil {
		log.Fatalf("Invalid URL: %s", err)
	}

	htmlContent, err := source.FetchContent(validURL)
	if err != nil {
		log.Fatalf("Failed to fetch: %s", err)
	}

	if err := source.AnalyzeEntities(htmlContent); err != nil {
		log.Fatalf("Failed to analyse entities: %s", err)
	}

	// Get SERP results for a given keyword
	keyword, err := source.KwInput()
	if err != nil {
		log.Fatalf("Invalid keyword: %s", err)
	}

	validKeyword, err := source.ValidateKeyword(keyword)
	if err != nil {
		log.Fatalf("Unable to process the target keyword")
	}

	serpresult, err := source.SerpExtraction(validKeyword)
	if err != nil {
		log.Fatalf("Couldn't fetch for %s", keyword)
	}
	fmt.Println("Top Google Search results:")
	for i, url := range serpresult[2:] {
		fmt.Printf("%d: %s \n", i+1, url)
	}*/

	entities := []source.Entity{
		{Name: "trans", WikiURLfromWiki: "https://en.wikipedia.org/wiki/Trans"},
		{Name: "transgender", WikiURLfromWiki: "https://en.wikipedia.org/wiki/Transgender"},
		{Name: "cd", WikiURLfromWiki: "https://en.wikipedia.org/wiki/Cross-dressing"},
	}

	schema, err := source.GenerateSchema(entities)
	if err != nil {
		log.Fatalf("Wansn't able to generate schema: %s ", err)
	}
	fmt.Println(schema)

}
