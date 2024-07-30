package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ProceedSchema() string {
	fmt.Println("Would you like to generate schema markups? - 'Yes' or any key to leave")
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
				"content": "you are an expert web developer and simply share the code",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Generate semantic schema markup for the following entities with their Wikipedia URLs:\n%s", entitiesSchema),
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

	// Unescape JSON string
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "    "); err != nil {
		return "", fmt.Errorf("error formatting JSON: %w", err)
	}

	return prettyJSON.String(), nil
}
