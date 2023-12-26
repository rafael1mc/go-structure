package auth

import (
	"fmt"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/jwt"
	"gomodel/internal/shared/redis"
	"gomodel/internal/shared/timeprovider"
	"gomodel/internal/shared/util/password"
	"gomodel/internal/user"
	"time"
)

type Auth struct {
	user         *user.User
	jwtWrapper   *jwt.JwtWrapper
	env          *env.Env
	timeProvider timeprovider.TimeProvider
	redisWrapper *redis.RedisWrapper
}

func NewAuth(
	user *user.User,
	jwtWrapper *jwt.JwtWrapper,
	env *env.Env,
	timeProvider timeprovider.TimeProvider,
	redisWrapper *redis.RedisWrapper,
) *Auth {
	return &Auth{
		user:         user,
		jwtWrapper:   jwtWrapper,
		env:          env,
		timeProvider: timeProvider,
		redisWrapper: redisWrapper,
	}
}

func (a *Auth) AuthenticateUser(authReq AuthRequest) (*AuthResponse, error) {
	dbUser, err := a.user.GetByEmail(authReq.Email, nil)
	if err != nil {
		return nil, err
	}

	if !dbUser.IsEnabled {
		return nil, UserNotEnabledError
	}

	if err := password.ComparedPwd(authReq.PlainPass, dbUser.PassHash, dbUser.PassSalt); err != nil {
		return nil, InvalidPasswordError
	}

	return a.createAuthResponse(dbUser.ID)
}

func (a *Auth) ValidateAccessToken(accessToken string) (userID string, err error) {
	tokenClaims, err := a.jwtWrapper.ValidateToken(accessToken)
	if err != nil {
		return
	}

	accessTokenKey := fmt.Sprintf("access:%s", tokenClaims.TokenUUID)
	userID, err = a.redisWrapper.Get(accessTokenKey)
	if err != nil {
		return
	}

	if userID == "" {
		err = fmt.Errorf("unauthorized")
		return
	}

	return
}

func (a *Auth) RefreshAccessToken(refreshTokenReq RefreshTokenRequest) (*AuthResponse, error) {
	if refreshTokenReq.RefreshToken == "" {
		return nil, fmt.Errorf("invalid refresh token")
	}

	tokenClaims, err := a.jwtWrapper.ValidateToken(refreshTokenReq.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenKey := fmt.Sprintf("refresh:%s", tokenClaims.TokenUUID)
	userID, err := a.redisWrapper.Get(refreshTokenKey)
	if err != nil {
		return nil, err
	}

	dbUser, err := a.user.GetByID(userID, nil)
	if err != nil {
		return nil, err
	}

	if !dbUser.IsEnabled {
		return nil, fmt.Errorf("unauthorized")
	}

	return a.createAuthResponse(dbUser.ID)
}

func (a *Auth) Logout(userID string) (err error) {
	redisKeys := []string{}
	err = a.redisWrapper.GetWithPatternIterable("access:*", func(key, value string) bool {
		if value == userID {
			redisKeys = append(redisKeys, key)
		}
		return false
	})
	if err != nil {
		return
	}
	err = a.redisWrapper.GetWithPatternIterable("refresh:*", func(key, value string) bool {
		if value == userID {
			redisKeys = append(redisKeys, key)
		}
		return false
	})
	if err != nil {
		return
	}

	a.redisWrapper.Del(redisKeys...)

	return nil
}

func (a *Auth) createAuthResponse(userID string) (*AuthResponse, error) {
	accessTokenDuration := time.Duration(a.env.Api.AccessTokenTTLSeconds) * time.Second
	accessToken, err := a.jwtWrapper.CreateToken(userID, accessTokenDuration, a.env.Api.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	refreshTokenDuration := time.Duration(a.env.Api.RefreshTokenTTLSeconds) * time.Second
	refreshToken, err := a.jwtWrapper.CreateToken(userID, refreshTokenDuration, a.env.Api.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	accessTokenKey := fmt.Sprintf("access:%s", accessToken.TokenUUID)
	err = a.redisWrapper.Save(accessTokenKey, userID, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshTokenKey := fmt.Sprintf("refresh:%s", refreshToken.TokenUUID)
	err = a.redisWrapper.Save(refreshTokenKey, userID, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
