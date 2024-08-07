package main

import (
	"fmt"
	"nlp-tool/cmd"
	"nlp-tool/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("You won't able to generate schema markups: %s", err)
	}
	cmd.Execute()
}
