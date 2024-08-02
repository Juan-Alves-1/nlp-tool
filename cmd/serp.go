package cmd

import (
	"fmt"
	"log"
	"nlp-tool/internal"

	"github.com/spf13/cobra"
)

var serpCmd = &cobra.Command{
	Use:   "serp",
	Short: "Extract SERP data for a given keyword on Google Search US",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cm *cobra.Command, args []string) {
		var keyword string
		if len(args) > 0 {
			keyword = args[0]
		} else {
			var err error
			keyword, err = internal.KwInput()
			if err != nil {
				log.Fatalf("Invalid keyword: %s", err)
			}
		}

		validKeyword, err := internal.ValidateKeyword(keyword)
		if err != nil {
			log.Fatalf("Unable to process the target keyword")
		}

		serpresult, err := internal.SerpExtraction(validKeyword)
		if err != nil {
			log.Fatalf("Couldn't fetch for %s", keyword)
		}
		fmt.Println(" ")
		fmt.Println("Top Google Search results:")
		for i, url := range serpresult[2:] {
			fmt.Printf("%d: %s \n", i+1, url)
		}
	},
}

func init() {
	rootCmd.AddCommand(serpCmd)
}
