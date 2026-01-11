package driver

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

func NewGenAIClient() *genai.Client {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if googleAPIKey == "" && geminiAPIKey == "" {
		log.Fatal("genai: required environment variable (GOOGLE_API_KEY or GEMINI_API_KEY) is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("genai: failed to create client: %v", err)
	}

	return client
}
