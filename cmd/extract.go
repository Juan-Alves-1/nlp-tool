package cmd

import (
	"fmt"
	"nlp-tool/internal"
	"os"

	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract top entities from a given URL",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		var err error
		if len(args) > 0 {
			url = args[0]
		} else {
			url, err = internal.UrlInput()
			if err != nil {
				fmt.Printf("Invalid input: %s\n", err)
				os.Exit(2)
			}
		}

		validURL, err := internal.ValidateURL(url)
		if err != nil {
			fmt.Printf("Invalid URL: %s\n", err)
			os.Exit(3)
		}

		htmlContent, err := internal.FetchContent(validURL)
		if err != nil {
			fmt.Printf("Failed to fetch: %s\n", err)
			os.Exit(4)
		}

		topEntities, err := internal.AnalyzeEntities(htmlContent)
		if err != nil {
			fmt.Printf("Failed to analyse entities: %s\n", err)
			os.Exit(5)
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
			fmt.Printf("Wansn't able to generate schema: %s ", err)
			os.Exit(6)
		}

		fmt.Println("\n Generated schema:\n", schema)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
