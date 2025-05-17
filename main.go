package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/suyashmohan/go-aiagent/cmd"
	"github.com/suyashmohan/go-aiagent/internal"
	"github.com/suyashmohan/go-aiagent/internal/tools"
	"github.com/urfave/cli/v3"
)

func main() {
	// Load ENV from .env file
	godotenv.Load()

	// Load System prompt
	fileBytes, err := os.ReadFile("system.md")
	if err != nil {
		log.Fatalln(err)
	}
	systemPrompt := string(fileBytes)

	// Create an Agent with Tools
	agent, err := internal.NewAgentWithTools(systemPrompt, getTools())
	if err != nil {
		log.Fatalln("failed to create ai agent", err)
	}

	// Setup CLI commands
	cmd := &cli.Command{
		Name:   "agent",
		Usage:  "Run AI agent with a query",
		Action: cmd.RootCMD(agent),
		// Separate command to start as Slack bot
		Commands: []*cli.Command{
			{
				Name:   "slack",
				Usage:  "Run AI Agent as slack bot",
				Action: cmd.SlackCMD(agent),
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}
}

func getTools() map[string]internal.AgentTool {
	weatherToolName, weatherTool := tools.WeatherTool()
	return map[string]internal.AgentTool{
		weatherToolName: weatherTool,
	}
}
