package go_certcentral

type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	JobTitle    string `json:"job_title,omitempty"`
	Telephone   string `json:"telephone,omitempty"`
	Name        string `json:"name,omitempty"`
	ContactType string `json:"contact_type,omitempty"`
}
