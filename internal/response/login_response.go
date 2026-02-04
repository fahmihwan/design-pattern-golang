package response

import "time"

type LoginResponse struct {
	User        any       `json:"user"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	TokenType   string    `json:"token_type"`
}
