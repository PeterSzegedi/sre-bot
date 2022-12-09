package oncall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"skippy/pkg/slack/message"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetOnCallSRE(token string, schedule string, client *slack.Client, channel string, current bool, L string) {
	// Default to getting the current oncall SRE
	query := fmt.Sprintf("https://api.pagerduty.com/oncalls?include[]=users&schedule_ids[]=%s&earliest=true", schedule)

	// Get the SREs in the next 12 hours from Pagerduty
	if !current {
		since := time.Now().Format(time.RFC3339)
		until := time.Now().Add(time.Hour * 12).Format(time.RFC3339)
		query = fmt.Sprintf("https://api.pagerduty.com/oncalls?include[]=users&schedule_ids[]=%s&since=%s&until=%s", schedule, since, until)
	}

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", token))
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Content-Type", "application/json")

	hClient := &http.Client{}
	resp, err := hClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var parsedBody OncallResponse
	err = json.Unmarshal(body, &parsedBody)

	if err != nil {
		log.Fatalln(err)
	}
	log.Info(string(body))
	slackBlock := ""
	if len(parsedBody.Oncalls) > 0 {
		slackBlock = parsedBody.Oncalls[0].User.Name
	}
	if !current {
		if len(parsedBody.Oncalls) > 1 {
			slackBlock = parsedBody.Oncalls[1].User.Name
		}
	}
	if slackBlock != "" {
		message.PostMessage(client, channel, slackBlock, "", fmt.Sprintf("The current %s SRE", L))
	} else {
		message.PostMessage(client, channel, slackBlock, "", "Pagerduty API has returned an empty oncall list")
	}
}
