package cmd

import (
	"fmt"
	"nlp-tool/internal"
	"os"

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
				fmt.Printf("Invalid keyword: %s\n", err)
				os.Exit(7)

			}
		}

		validKeyword, err := internal.ValidateKeyword(keyword)
		if err != nil {
			fmt.Println("Unable to process the target keyword")
			os.Exit(8)
		}

		serpresult, err := internal.SerpExtraction(validKeyword)
		if err != nil {
			fmt.Printf("Couldn't fetch for %s\n", keyword)
			os.Exit(9)
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
