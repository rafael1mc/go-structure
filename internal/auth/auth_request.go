package auth

type AuthRequest struct {
	Email     string `json:"email"`
	PlainPass string `json:"password"`
}
