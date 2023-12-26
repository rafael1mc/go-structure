package jwt

import (
	"encoding/base64"
	"fmt"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/timeprovider"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtWrapper struct {
	timeProvider timeprovider.TimeProvider
	env          *env.Env
}

func NewJwtWrapper(
	timeProvider timeprovider.TimeProvider,
	env *env.Env,
) *JwtWrapper {
	return &JwtWrapper{
		timeProvider: timeProvider,
		env:          env,
	}
}

func (w *JwtWrapper) CreateToken(userID string, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := w.timeProvider.ProvideUTCSec()
	td := TokenDetails{
		ExpiresAt: new(int64),
		Token:     new(string),
	}

	*td.ExpiresAt = now + int64(ttl.Seconds())
	td.TokenUUID = uuid.NewString()
	td.UserID = userID

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(w.env.Api.AccessTokenPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	mapClaims := jwt.MapClaims{
		"sub":        userID,
		"token_uuid": td.TokenUUID,
		"exp":        td.ExpiresAt,
		"iat":        now,
		"nbf":        now,
	}

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	return &td, nil
}

func (w *JwtWrapper) ValidateToken(token string) (*TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(w.env.Api.AccessTokenPublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}
