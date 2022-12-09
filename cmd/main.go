package main

import (
	"skippy/pkg/env"
	"skippy/pkg/health"
	"skippy/pkg/slack/connect"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Get the tokens
	env := env.PopulateEnv()

	// Start the Slack connection
	api, client := connect.InitSlackConnection(env.BotToken, env.AppToken, env.Debug)

	// Test the connection and return the current users
	currentUser := connect.TestSlackConnection(api)

	// spin up the currently stub healthcheck
	health.InitHealth()
	// Run the infinite loop which will monitor the Slack events for bot mentions
	err := connect.RunEventLoop(api, client, env, currentUser)
	if err != nil {
		log.Fatalf("Failed run socketmode %s", err)
	}
}
