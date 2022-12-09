package message

import (
	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
)

// postMessage will send messages to the event's channel
func PostMessage(client *slack.Client, channel string, slackBlock string, threadTS string, text string) {
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
