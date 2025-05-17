package cmd

import (
	"context"
	"fmt"

	"github.com/suyashmohan/go-aiagent/internal"
	"github.com/urfave/cli/v3"
)

func SlackCMD(agent *internal.Agent) func(ctx context.Context, c *cli.Command) error {
	return func(ctx context.Context, c *cli.Command) error {
		slackBot := internal.NewSlackBot(agent)
		if slackBot.Run() != nil {
			return fmt.Errorf("failed to run slack agent ")
		}

		return nil
	}
}
