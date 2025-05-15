package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/suyashmohan/go-aiagent/internal"
	"github.com/urfave/cli/v3"
)

func RootCMD(agent *internal.Agent) func(ctx context.Context, c *cli.Command) error {
	return func(ctx context.Context, c *cli.Command) error {
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
	}
}
