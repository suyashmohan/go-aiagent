package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Agent struct {
	SystemPrompt string
	ModelName    string
	client       openai.Client
}

func NewAgent(systemPrompt string) (*Agent, error) {
	API_KEY := os.Getenv("OPENAI_API_KEY")
	BASE_URL := os.Getenv("API_BASE_URL")
	MODEL_NAME := os.Getenv("OPENAI_MODEL_NAME")

	if API_KEY == "" || BASE_URL == "" || MODEL_NAME == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY, BASE_URL and MODEL_NAME must be set in environment variables")
	}

	return &Agent{
		SystemPrompt: systemPrompt,
		ModelName:    MODEL_NAME,
		client: openai.NewClient(
			option.WithBaseURL(BASE_URL),
			option.WithAPIKey(API_KEY),
		),
	}, nil
}

func (a *Agent) Run(context context.Context, prompt string) (string, error) {
	chatCompletion, err := a.client.Chat.Completions.New(context, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(a.SystemPrompt),
			openai.UserMessage(prompt),
		},
		Model: a.ModelName,
	})

	if err != nil {
		return "", fmt.Errorf("failed to call openai chat completions api - %w", err)
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
