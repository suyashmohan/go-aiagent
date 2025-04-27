package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func StartSlackBot() error {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if appToken == "" || botToken == "" {
		return fmt.Errorf("SLACK_APP_TOKEN and SLACK_BOT_TOKEN must be set in environment variables")
	}

	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(
		api,
	)
	socketModeHandler := socketmode.NewSocketmodeHandler(
		client,
	)
	socketModeHandler.HandleEvents(slackevents.AppMention, handleAppMentionEvent)
	socketModeHandler.RunEventLoop()

	return nil
}

func handleAppMentionEvent(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("ignored %+v\n", evt)
	}
	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		log.Printf("ignored %+v\n", ev)
	}

	_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Hello", false))
	if err != nil {
		log.Printf("failed to post message on slack - %s", err)
	}
}
