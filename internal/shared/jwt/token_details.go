package jwt

type TokenDetails struct {
	Token     *string `json:"token"`
	TokenUUID string  `json:"token_uuid"`
	UserID    string  `json:"user_id"`
	ExpiresAt *int64  `json:"expires_at"`
}
