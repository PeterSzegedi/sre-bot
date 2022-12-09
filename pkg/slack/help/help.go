package help

import (
	"fmt"
	"skippy/pkg/slack/message"

	"github.com/slack-go/slack"
)

// helpText sends help text for the caller to Slack
func HelpText(client *slack.Client, firstName string, channel string, timeStamp string) {
	text := fmt.Sprintf("Hi %s", firstName)
	slackMessage := `Currently I have the following functions:
	- standup-order
	- find-current-l1
	- find-next-l1
	- find-current-l2
	- find-next-l2`
	message.PostMessage(client, channel, slackMessage, timeStamp, text)
}
