package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func ProceedSchema() string {
	fmt.Println(" ")
	fmt.Println("Would you like to generate schema markups?")
	fmt.Println(" Yes to proceed or any key to leave")
	var input string
	fmt.Scanln(&input)
	if input == "yes" || input == "y" || input == "Yes" {
		return "Analysing your entities..."
	}
	return "Bye! Have a good one :)"
}

func GenerateSchema(entities []Entity) (string, error) {
	// Load .env file
	err := godotenv.Load("../nlp-tool/.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return "", fmt.Errorf("error loading .env file: %v", err)
	}
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Construct the entities string for the prompt
	var entitiesSchema string
	for _, entity := range entities {
		entitiesSchema += fmt.Sprintf("Entity Name: %s, Wikipedia URL: %s\n", entity.Name, entity.WikiURLfromWiki)
	}

	// Define the request payload
	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "you are an expert web developer in schema markups. Use this snippet as a reference for requests: <script type='application/ld+json'> {'@context': 'https://schema.org','@type': 'WebPage', 'mentions': [{'@type': 'Organization','name': 'Samsung','sameAs': 'https://en.wikipedia.org/wiki/Samsung'},{'@type': 'Thing','name': 'OLED','sameAs': 'https://en.wikipedia.org/wiki/OLED'},{'@type': 'Thing','name': 'selection','sameAs': 'https://en.wikipedia.org/wiki/Choice'},{'@type': 'Thing','name': 'picture quality','sameAs': 'https://en.wikipedia.org/wiki/Image'},{'@type': 'Thing','name': 'market','sameAs': 'https://en.wikipedia.org/wiki/Market_(economics)'},{'@type': 'Thing','name': 'price','sameAs': 'https://en.wikipedia.org/wiki/Price'},{'@type': 'Thing','name': 'High-end','sameAs': 'https://en.wikipedia.org/wiki/Luxury_goods'}]}</script> ",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("only provide the code as an output: generate schema markups similar to the your reference and note that, from mentions onwards, match 'name' keys with the following Entity Names as values, and for 'sameAs' keys matche their following Wikipedia URLs values:\n%s", entitiesSchema),
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %s, body: %s", resp.Status, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Extract and format the schema content
	schemaContent := openAIResp.Choices[0].Message.Content

	return schemaContent, nil
}
