package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
}

var AppConfig Config

func LoadConfig() error {
	err := godotenv.Load("../nlp-tool/.env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	AppConfig = Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}

	if AppConfig.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set in the environment variables")
	}

	return nil
}
