package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
		Name:  "agent",
		Usage: "Run AI agent with a query",
		Action: func(ctx context.Context, c *cli.Command) error {
			// Main cli command that call agent with cli param
			if c.Args().Len() == 0 {
				log.Println("No query provided.")
				return nil
			}
			query := c.Args().First()
			answer, err := agent.Run(context.TODO(), query)
			if err != nil {
				return fmt.Errorf("failed to run ai agent - %w", err)
			}
			fmt.Println(answer)

			return nil
		},
		// Separate command to start as Slack bot
		Commands: []*cli.Command{
			{
				Name:  "slack",
				Usage: "Run AI Agent as slack bot",
				Action: func(ctx context.Context, c *cli.Command) error {
					slackBot := internal.NewSlackBot(agent)
					if slackBot.Run() != nil {
						return fmt.Errorf("failed to run slack agent - %w", err)
					}

					return nil
				},
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
