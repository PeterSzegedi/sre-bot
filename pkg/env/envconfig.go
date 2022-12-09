package env

// EnvConfig stores the Slack tokens
type EnvConfig struct {
	// AppToken is app-level-token to run socketmode.
	AppToken string `envconfig:"APP_LEVEL_TOKEN" required:"true"`
	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
	// PDToken is a read-only token to access Pagerduty schedules.
	PDToken string `envconfig:"PD_TOKEN" required:"true"`
	// L1Schedule is the ID of the L1 schedule to be retrieved from Pagerduty by the bot
	L1Schedule string `envconfig:"L1_SCHEDULE" required:"true"`
	// L2Schedule is the ID of the L2 schedule to be retrieved from Pagerduty by the bot
	L2Schedule string `envconfig:"L2_SCHEDULE" required:"true"`
	// Enabling debug mode. This is a stub for now
	Debug bool `envconfig:"DEBUG" required:"false" default:"false"`
}
