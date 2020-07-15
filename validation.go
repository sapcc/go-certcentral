package go_certcentral

import "time"

type Validation struct {
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	DateCreated    *time.Time `json:"date_created,omitempty"`
	ValidatedUntil *time.Time `json:"validated_until,omitempty"`
	Status         string     `json:"status"`
	DcvStatus      string     `json:"dcv_status,omitempty"`
	VerifiedUsers  []User     `json:"verified_users,omitempty"`
}
