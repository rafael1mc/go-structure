package auth

import "gomodel/internal/shared/jwt"

type AuthResponse struct {
	AccessToken  *jwt.TokenDetails `json:"access_token"`
	RefreshToken *jwt.TokenDetails `json:"refresh_token"`
}
