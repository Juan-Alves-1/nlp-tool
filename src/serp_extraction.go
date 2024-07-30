package source

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func ProceedSerpExtraction() string {
	fmt.Println("Would you like to see the top SERPs for a given keyword?")
	fmt.Println("Yes to proceed or any key to leave")
	var input string
	fmt.Scanln(&input)
	if input == "yes" || input == "y" || input == "Yes" {
		return "Fetching Google Search US..."
	}
	return "Bye! Have a good one :)"
}

func SerpExtraction(keyword string) ([]string, error) {
	collector := colly.NewCollector()

	var topResultURLs []string

	// Find and visit the first 15 links in the search results
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/url?q=") && len(topResultURLs) < 17 {
			url := strings.Split(href, "&")[0][7:] // Remove "/url?q=" and everything after "&"
			topResultURLs = append(topResultURLs, url)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("___________________")
		fmt.Println("Visiting", r.URL)
	})

	cleanPage := "https://www.google.com/search?q=" + keyword + "&hl=en&gl=us"

	err := collector.Visit(cleanPage)
	if err != nil {
		log.Fatalf("Failed to visit Google search page: %v", err)
	}

	return topResultURLs, nil

	/*fmt.Println("Top 10 Google result URLs:")
	for i, url := range topResultURLs[2:] {
		fmt.Printf("%d: %s\n", i+1, url)
	}*/
}
