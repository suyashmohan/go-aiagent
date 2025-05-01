package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/suyashmohan/myaiagent/internal"
	"github.com/urfave/cli/v3"
)

func main() {
	godotenv.Load()
	fileBytes, err := os.ReadFile("system.md")
	if err != nil {
		log.Fatalln(err)
	}
	systemPrompt := string(fileBytes)
	cmd := &cli.Command{
		Name:  "agent",
		Usage: "Run AI agent with a query",
		Action: func(ctx context.Context, c *cli.Command) error {
			agent, err := internal.NewAgent(systemPrompt)
			if err != nil {
				return fmt.Errorf("failed to create ai agent - %w", err)
			}

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
		Commands: []*cli.Command{
			{
				Name:  "slack",
				Usage: "Run AI Agent as slack bot",
				Action: func(ctx context.Context, c *cli.Command) error {
					agent, err := internal.NewAgent(systemPrompt)
					if err != nil {
						return fmt.Errorf("failed to create ai agent - %w", err)
					}

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
