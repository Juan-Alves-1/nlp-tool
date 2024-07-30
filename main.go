package main

import (
	"fmt"
	"log"
	source "nlp-tool/src"
)

func main() {
	// Get top entities for a given URL
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

	topEntities, err := source.AnalyzeEntities(htmlContent)
	if err != nil {
		log.Fatalf("Failed to analyse entities: %s", err)
	}

	for _, entity := range topEntities {
		fmt.Printf("Name: %s, Salience: %.3f, Type: %s\n",
			entity.Name, entity.Salience, entity.Type)
	}

	// Generate schema markups based on previous entities
	schemaMessage := source.ProceedSchema()
	if schemaMessage != "Analysing your entities..." {
		fmt.Println(schemaMessage)
		return
	}
	fmt.Println(schemaMessage)

	schema, err := source.GenerateSchema(topEntities[:5])
	if err != nil {
		log.Fatalf("Wansn't able to generate schema: %s ", err)
	}
	fmt.Println("Generated schema:\n", schema)

	serpMessage := source.ProceedSerpExtraction()
	if serpMessage != "Fetching Google Search US..." {
		fmt.Println(serpMessage)
	}
	fmt.Println(serpMessage)

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
	}
}
