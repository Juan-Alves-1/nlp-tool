package cmd

import (
	"fmt"
	"log"
	"nlp-tool/internal"

	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract top entities from a given URL",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) > 0 {
			url = args[0]
		} else {
			var err error
			url, err = internal.UrlInput()
			if err != nil {
				log.Fatalf("Invalid input: %s", err)
			}
		}

		validURL, err := internal.ValidateURL(url)
		if err != nil {
			log.Fatalf("Invalid URL: %s", err)
		}

		htmlContent, err := internal.FetchContent(validURL)
		if err != nil {
			log.Fatalf("Failed to fetch: %s", err)
		}

		topEntities, err := internal.AnalyzeEntities(htmlContent)
		if err != nil {
			log.Fatalf("Failed to analyse entities: %s", err)
		}

		for _, entity := range topEntities {
			fmt.Printf("Name: %s, Salience: %.3f, Type: %s\n",
				entity.Name, entity.Salience, entity.Type)
		}

		// Generate schema markups based on previous entities
		schemaMessage := internal.ProceedSchema()
		if schemaMessage != "Analysing your entities..." {
			fmt.Println(schemaMessage)
			return
		}
		fmt.Println(schemaMessage)

		schema, err := internal.GenerateSchema(topEntities[:10])
		if err != nil {
			log.Fatalf("Wansn't able to generate schema: %s ", err)
		}
		fmt.Println(" ")
		fmt.Println("Generated schema:\n", schema)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
