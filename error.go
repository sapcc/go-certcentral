package go_certcentral

import "fmt"

type Error struct {
	Code    int    `json:"-"`
	Status  string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, status: %s, message: %s", e.Code, e.Status, e.Message)
}
