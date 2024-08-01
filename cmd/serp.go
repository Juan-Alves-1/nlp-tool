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
	Run: func(cm *cobra.Command, arg []string) {
		serpMessage := internal.ProceedSerpExtraction()
		if serpMessage != "Sure!" {
			fmt.Println(serpMessage)
			return
		}
		fmt.Println(serpMessage)

		keyword, err := internal.KwInput()
		if err != nil {
			log.Fatalf("Invalid keyword: %s", err)
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
