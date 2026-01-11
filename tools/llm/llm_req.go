package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hackathon-20260110/api/adapter"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("failed to create genai client: %v", err)
	}

	llmAdapter := adapter.NewLLMAdapter(client)
	resp, err := llmAdapter.CreateChatCompletion([]*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: "Hello, how are you?",
				},
			},
		},
	}, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		log.Fatalf("CreateChatCompletion error: %v", err)
	}
	fmt.Println(resp)
}
