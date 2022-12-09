package unknown

import (
	"fmt"
	"skippy/pkg/slack/message"

	"github.com/slack-go/slack"
)

// unkownFunction is a catch all function for texts that are unknown to the bot (as functions)
func UnknownFunction(client *slack.Client, firstName string, channel string, timeStamp string) {
	text := fmt.Sprintf("I'm Afraid I Can't Do That, %s", firstName)
	slackMessage := "Currently I have the following functions:\n- standup-order"
	message.PostMessage(client, channel, slackMessage, timeStamp, text)
}
