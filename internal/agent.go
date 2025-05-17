package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Agent struct {
	SystemPrompt string
	ModelName    string
	Tools        map[string]AgentTool
	client       openai.Client
}

type AgentTool struct {
	Description string
	Parameters  map[string]any
	Required    []string
	Fn          func(map[string]interface{}) string
}

func NewAgent(systemPrompt string) (*Agent, error) {
	return NewAgentWithTools(systemPrompt, map[string]AgentTool{})
}

func NewAgentWithTools(systemPrompt string, tools map[string]AgentTool) (*Agent, error) {
	API_KEY := os.Getenv("OPENAI_API_KEY")
	BASE_URL := os.Getenv("API_BASE_URL")
	MODEL_NAME := os.Getenv("OPENAI_MODEL_NAME")

	if API_KEY == "" || BASE_URL == "" || MODEL_NAME == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY, BASE_URL and MODEL_NAME must be set in environment variables")
	}

	return &Agent{
		SystemPrompt: systemPrompt,
		ModelName:    MODEL_NAME,
		Tools:        tools,
		client: openai.NewClient(
			option.WithBaseURL(BASE_URL),
			option.WithAPIKey(API_KEY),
		),
	}, nil
}

func (a *Agent) Run(context context.Context, prompt string) (string, error) {
	maxSteps := 5
	currStep := 0

	// Prepare Tools for function calling
	toolFuncs := []openai.ChatCompletionToolParam{}
	for name, tool := range a.Tools {
		// Convert our style of Function definition into OpenAI style
		toolFuncs = append(toolFuncs, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        name,
				Description: openai.String(tool.Description),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"required":   tool.Required,
					"properties": tool.Parameters,
				},
			},
		})
	}

	// Prepare Messages
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(a.SystemPrompt),
		openai.UserMessage(prompt),
	}

	// Iterate over Chat Completions API based on tool usage
	for currStep < maxSteps {
		logJson, _ := json.MarshalIndent(messages, "", "  ")
		log.Println("Step", currStep)
		log.Println(string(logJson))

		// Call Chat Completions API
		chatCompletion, err := a.client.Chat.Completions.New(context, openai.ChatCompletionNewParams{
			Messages: messages,
			Tools:    toolFuncs,
			Model:    a.ModelName,
		})

		if err != nil {
			return "", fmt.Errorf("failed to call openai chat completions api - %w", err)
		}

		// Check if there is a need to call a tool
		toolCalls := chatCompletion.Choices[0].Message.ToolCalls
		if len(toolCalls) > 0 {
			messages = append(messages, chatCompletion.Choices[0].Message.ToParam())
			for _, toolCall := range toolCalls {
				log.Println("Request to call", toolCall.Function.Name)
				// Extract the args for tool calling
				var args map[string]interface{}
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				if err != nil {
					return "", fmt.Errorf("failed to extract function arguments - %w", err)
				}
				answer := a.Tools[toolCall.Function.Name].Fn(args)
				messages = append(messages, openai.ToolMessage(answer, toolCall.ID))
			}
		} else {
			return chatCompletion.Choices[0].Message.Content, nil
		}

		currStep = currStep + 1
	}

	return "", fmt.Errorf("something failed when calling openai chat completions api")
}
