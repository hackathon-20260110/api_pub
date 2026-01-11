package adapter

import (
	"context"

	"github.com/hackathon-20260110/api/utils"
	"google.golang.org/genai"
)

type LLMModelType string

const (
	LLM_MODEL_TYPE_GEMINI2_0             LLMModelType = "gemini-2.0-flash"
	LLM_MODEL_TYPE_GEMINI2_5_FLASH       LLMModelType = "gemini-2.5-flash"
	LLM_MODEL_TYPE_GEMINI_2_5_FLASH_LITE LLMModelType = "gemini-2.5-flash-lite-preview-09-2025"
)

type LLMAdapter interface {
	CreateChatCompletion(contents []*genai.Content, model LLMModelType) (string, error)
	CreateChatCompletionJSON(contents []*genai.Content, model LLMModelType) (string, error)
}

func NewLLMAdapter(genaiClient *genai.Client) LLMAdapter {
	return &llmAdapter{client: genaiClient}
}

type llmAdapter struct {
	client *genai.Client
}

func (a *llmAdapter) CreateChatCompletion(contents []*genai.Content, model LLMModelType) (string, error) {
	ctx := context.Background()

	result, err := a.client.Models.GenerateContent(
		ctx,
		string(model),
		contents,
		nil,
	)
	if err != nil {
		return "", utils.WrapError(err)
	}

	return result.Text(), nil
}

func (a *llmAdapter) CreateChatCompletionJSON(contents []*genai.Content, model LLMModelType) (string, error) {
	ctx := context.Background()

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	result, err := a.client.Models.GenerateContent(
		ctx,
		string(model),
		contents,
		config,
	)
	if err != nil {
		return "", utils.WrapError(err)
	}

	return result.Text(), nil
}
