package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nlp-tool",
	Short: "NLP Tool is a CLI application for entity extraction and SERP data mining",
}

func Execute() {
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(serpCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
