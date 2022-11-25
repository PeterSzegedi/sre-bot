package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var HEALTHY = true

// EnvConfig stores the Slack tokens
type EnvConfig struct {
	// AppToken is app-level-token to run socketmode.
	AppToken string `envconfig:"APP_LEVEL_TOKEN" required:"true"`
	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
	//
	Debug bool `envconfig:"DEBUG" required:"false" default:"false"`
}

// Member to store Slack userdata
type Member struct {
	ID   string
	Name string
}

// Populate tokens from environment variables
func populateEnv() EnvConfig {
	var env EnvConfig
	err := envconfig.Process("", &env)
	if err != nil {
		log.Error("Failed to read env vars: %s", err)
	}
	return env
}

// initSlackConnection initiates a socketmode connection to the Slack API
func initSlackConnection(botToken string, appToken string, debugMode bool) (*slack.Client, *socketmode.Client) {
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
func testSlackConnection(client *slack.Client) string {
	authTest, authTestErr := client.AuthTest()
	if authTestErr != nil {
		log.Fatalf("SLACK_BOT_TOKEN is invalid: %v", authTestErr)
		os.Exit(1)
	}
	return authTest.UserID
}

// runEventLoop monitors the channel for mentions of the sre-bot
func runEventLoop(client *slack.Client, socketMode *socketmode.Client, env EnvConfig, currentUser string) error {
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
						helpText(client, userInfo.Profile.FirstName, ev.Channel, ev.TimeStamp)

					} else if strings.Split(ev.Text, " ")[1] == "standup-order" {
						// bot mention for standup-order, print randomized members of the channel
						standupOrder(client, ev.Channel)

					} else {
						// bot mention for unknown function, print help
						unknownFunction(client, userInfo.Profile.FirstName, ev.Channel, ev.TimeStamp)
					}

				}
			}
		}
	}()

	return socketMode.Run()
}

// helpText sends help text for the caller to Slack
func helpText(client *slack.Client, firstName string, channel string, timeStamp string) {
	text := fmt.Sprintf("Hi %s", firstName)
	message := "Currently I have the following functions:\n- standup-order"
	postMessage(client, channel, message, timeStamp, text)
}

// standupOrder generates a randomly selected numbered list of the members of the event's channel
func standupOrder(client *slack.Client, channel string) {
	allMembers, _ := fetchConversationMembers(channel, client)
	activeMembers, _ := filterActiveMembers(allMembers, client)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(activeMembers), func(i, j int) { activeMembers[i], activeMembers[j] = activeMembers[j], activeMembers[i] })
	log.Info(activeMembers)
	slackBlock := ""
	for i, member := range activeMembers {
		slackBlock += fmt.Sprintf("%v. %s\n", i+1, member.Name)

	}
	postMessage(client, channel, slackBlock, "", "Shuffling members of the channel")
}

// fetchConversationMembers collects all users in the event's channel
func fetchConversationMembers(channelID string, client *slack.Client) ([]string, error) {
	nextCursor := ""
	params := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Cursor:    nextCursor,
	}
	var members []string
	fetchedMembers, nextCursor, err := client.GetUsersInConversation(&params)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	members = append(members, fetchedMembers...)
	for len(nextCursor) > 0 {
		fetchedMembers, nextCursor, err = client.GetUsersInConversation(&params)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		members = append(members, fetchedMembers...)
	}
	return members, err
}

// filterActiveMember filters the users to non-bot and non-deleted users
func filterActiveMembers(members []string, client *slack.Client) ([]Member, error) {
	var activeMembers []Member
	for _, mem := range members {
		user, err := client.GetUserInfo(mem)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if !(user.IsBot || user.Deleted) {
			activeMembers = append(activeMembers, Member{ID: user.ID, Name: user.Name})
		}
	}
	return activeMembers, nil
}

// unkownFunction is a catch all function for texts that are unknown to the bot (as functions)
func unknownFunction(client *slack.Client, firstName string, channel string, timeStamp string) {
	text := fmt.Sprintf("I'm Afraid I Can't Do That, %s", firstName)
	message := "Currently I have the following functions:\n- standup-order"
	postMessage(client, channel, message, timeStamp, text)
}

// postMessage will send messages to the event's channel
func postMessage(client *slack.Client, channel string, slackBlock string, threadTS string, text string) {
	attachment := slack.Attachment{
		Text: slackBlock,
	}

	channelID, timestamp, err := client.PostMessage(
		channel,
		slack.MsgOptionText(text, false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionTS(threadTS),
	)
	if err != nil {
		log.Errorf("%s\n", err)
		return
	}
	log.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	body := json.RawMessage("OK")
	if !HEALTHY {
		body = json.RawMessage("NOT OK")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err := w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error writing response body: %s", err)
	}
}

func main() {
	// Get the tokens
	env := populateEnv()

	// Start the Slack connection
	api, client := initSlackConnection(env.BotToken, env.AppToken, env.Debug)

	// Test the connection and return the current users
	currentUser := testSlackConnection(api)

	healthServer := http.NewServeMux()
	healthServer.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r)
	})

	healthServer.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r)
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 8081), healthServer); err != nil {
			panic(err)
		}
	}()

	// Run the infinite loop which will monitor the Slack events for bot mentions
	err := runEventLoop(api, client, env, currentUser)
	if err != nil {
		log.Errorf("Failed run socketmode %s", err)
	}
}
