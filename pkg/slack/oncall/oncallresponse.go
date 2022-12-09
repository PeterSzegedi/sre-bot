package oncall

import "time"

type OncallResponse struct {
	Oncalls []struct {
		EscalationPolicy struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HTMLURL string `json:"html_url"`
		} `json:"escalation_policy"`
		EscalationLevel int `json:"escalation_level"`
		Schedule        struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HTMLURL string `json:"html_url"`
		} `json:"schedule"`
		User struct {
			Name           string      `json:"name"`
			Email          string      `json:"email"`
			TimeZone       string      `json:"time_zone"`
			Color          string      `json:"color"`
			AvatarURL      string      `json:"avatar_url"`
			Billed         bool        `json:"billed"`
			Role           string      `json:"role"`
			Description    interface{} `json:"description"`
			InvitationSent bool        `json:"invitation_sent"`
			JobTitle       string      `json:"job_title"`
			Teams          []struct {
				ID      string `json:"id"`
				Type    string `json:"type"`
				Summary string `json:"summary"`
				Self    string `json:"self"`
				HTMLURL string `json:"html_url"`
			} `json:"teams"`
			ContactMethods []struct {
				ID      string      `json:"id"`
				Type    string      `json:"type"`
				Summary string      `json:"summary"`
				Self    string      `json:"self"`
				HTMLURL interface{} `json:"html_url"`
			} `json:"contact_methods"`
			NotificationRules []struct {
				ID      string      `json:"id"`
				Type    string      `json:"type"`
				Summary string      `json:"summary"`
				Self    string      `json:"self"`
				HTMLURL interface{} `json:"html_url"`
			} `json:"notification_rules"`
			CoordinatedIncidents []interface{} `json:"coordinated_incidents"`
			ID                   string        `json:"id"`
			Type                 string        `json:"type"`
			Summary              string        `json:"summary"`
			Self                 string        `json:"self"`
			HTMLURL              string        `json:"html_url"`
		} `json:"user"`
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"oncalls"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	More   bool        `json:"more"`
	Total  interface{} `json:"total"`
}
