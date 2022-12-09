package connect

import (
	"skippy/pkg/env"
	"skippy/pkg/slack/help"
	"skippy/pkg/slack/oncall"
	"skippy/pkg/slack/standup"
	"skippy/pkg/slack/unknown"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// initSlackConnection initiates a socketmode connection to the Slack API
func InitSlackConnection(botToken string, appToken string, debugMode bool) (*slack.Client, *socketmode.Client) {
	api := slack.New(
		botToken,
		slack.OptionDebug(debugMode),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
	)
	return api, client
}

// testSlackConnection tests, if the connection was successful
func TestSlackConnection(client *slack.Client) string {
	authTest, authTestErr := client.AuthTest()
	if authTestErr != nil {
		log.Fatalf("SLACK_BOT_TOKEN is invalid: %v", authTestErr)
	}
	return authTest.UserID
}

// runEventLoop monitors the channel for mentions of the sre-bot
func RunEventLoop(client *slack.Client, socketMode *socketmode.Client, env env.EnvConfig, currentUser string) error {
	go func() {
		// loop to process events as they come
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				eventPayload := envelope.Data.(slackevents.EventsAPIEvent)
				switch ev := eventPayload.InnerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					socketMode.Ack(*envelope.Request)

					userInfo, _ := client.GetUserInfo(ev.User)
					if len(strings.Split(strings.Trim(ev.Text, " "), " ")) == 1 {
						// bot was mentioned without extra text, print help
						help.HelpText(client, userInfo.Profile.FirstName, ev.Channel, ev.TimeStamp)
					} else if strings.Split(ev.Text, " ")[1] == "standup-order" {
						// bot mention for standup-order, print randomized members of the channel
						standup.StandupOrder(client, ev.Channel)
					} else if strings.Split(ev.Text, " ")[1] == "find-current-l1" {
						// bot mention for find-current-l1, which finds the current SRE set on the L1 schedule on Pagerduty
						oncall.GetOnCallSRE(env.PDToken, env.L1Schedule, client, ev.Channel, true, "L1")
					} else if strings.Split(ev.Text, " ")[1] == "find-next-l1" {
						// bot mention for find-next-l1, which finds the next SRE set on the L1 schedule on Pagerduty
						oncall.GetOnCallSRE(env.PDToken, env.L1Schedule, client, ev.Channel, false, "L1")
					} else if strings.Split(ev.Text, " ")[1] == "find-current-l2" {
						// bot mention for find-current-l2, which finds the current SRE set on the L2 schedule on Pagerduty
						oncall.GetOnCallSRE(env.PDToken, env.L2Schedule, client, ev.Channel, true, "L2")
					} else if strings.Split(ev.Text, " ")[1] == "find-next-l2" {
						// bot mention for find-next-l2, which finds the next SRE set on the L2 schedule on Pagerduty
						oncall.GetOnCallSRE(env.PDToken, env.L2Schedule, client, ev.Channel, false, "L2")
					} else {
						// bot mention for unknown function, print help
						unknown.UnknownFunction(client, userInfo.Profile.FirstName, ev.Channel, ev.TimeStamp)
					}
				}
			}
		}
	}()

	return socketMode.Run()
}
