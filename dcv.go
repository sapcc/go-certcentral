package go_certcentral

import "time"

type (
	DCV struct {
		Method         string           `json:"method,omitempty"`
		NameScope      string           `json:"name_scope,omitempty"`
		DcvInvitations []DCVInvitations `json:"dcv_invitations,omitempty"`
	}

	DCVInvitations []struct {
		InvitationID int       `json:"invitation_id"`
		Email        string    `json:"email,omitempty"`
		Source       string    `json:"source,omitempty"`
		DateSent     time.Time `json:"date_sent,omitempty"`
		NameScope    string    `json:"name_scope,omitempty"`
	}
)
