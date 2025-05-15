package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackBot struct {
	agent *Agent
}

func NewSlackBot(agent *Agent) *SlackBot {
	return &SlackBot{
		agent,
	}
}

func (s *SlackBot) Run() error {
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
	socketModeHandler.HandleEvents(slackevents.AppMention, s.handleAppMentionEvent)
	socketModeHandler.RunEventLoop()

	return nil
}

func (s *SlackBot) handleAppMentionEvent(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("ignored %+v\n", evt)
	}
	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		log.Printf("ignored %+v\n", ev)
	}

	text := ev.Text
	re := regexp.MustCompile(`<@([A-Z0-9]+)>`)
	matches := re.FindAllStringSubmatch(text, -1)
	uniqeIDs := make(map[string]string)
	for _, match := range matches {
		if len(match) > 1 {
			uniqeIDs[match[1]] = match[0]
		}
	}
	for _, mention := range uniqeIDs {
		text = strings.ReplaceAll(text, mention, "")
	}
	log.Println("Received:", text)

	answer, aErr := s.agent.Run(context.Background(), text)
	if aErr != nil {
		log.Println("failed to run ai agent", aErr)
		answer = "failed to run ai agent"
	}
	log.Println("Answer:", answer)

	_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText(answer, false))
	if err != nil {
		log.Printf("failed to post message on slack - %s", err)
	}
}
