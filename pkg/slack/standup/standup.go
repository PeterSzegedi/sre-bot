package standup

import (
	"fmt"
	"math/rand"
	"skippy/pkg/slack/message"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
)

// Member to store Slack userdata
type Member struct {
	ID   string
	Name string
}

// standupOrder generates a randomly selected numbered list of the members of the event's channel
func StandupOrder(client *slack.Client, channel string) {
	allMembers, _ := fetchConversationMembers(channel, client)
	activeMembers, _ := filterActiveMembers(allMembers, client)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(activeMembers), func(i, j int) { activeMembers[i], activeMembers[j] = activeMembers[j], activeMembers[i] })
	log.Debug(activeMembers)
	slackBlock := ""
	for i, member := range activeMembers {
		slackBlock += fmt.Sprintf("%v. %s\n", i+1, member.Name)
	}
	message.PostMessage(client, channel, slackBlock, "", "Shuffling members of the channel")
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
