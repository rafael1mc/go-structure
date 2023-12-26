package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Environment Environment
	Api         Api
	Database    Database
	Redis       Redis
}

func NewEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v.\nContinuing...\n", err)
	}

	var env Env
	err = envconfig.Process("", &env)
	if err != nil {
		panic(err)
	}

	// Read access token private key
	bytes, err := os.ReadFile("keys/access.private")
	if err != nil {
		panic(err)
	}
	env.Api.AccessTokenPrivateKey = string(bytes)

	// Read access token public key
	bytes, err = os.ReadFile("keys/access.public")
	if err != nil {
		panic(err)
	}
	env.Api.AccessTokenPublicKey = string(bytes)

	// Read refresh token private key
	bytes, err = os.ReadFile("keys/refresh.private")
	if err != nil {
		panic(err)
	}
	env.Api.RefreshTokenPrivateKey = string(bytes)

	// Read refresh token public key
	bytes, err = os.ReadFile("keys/refresh.public")
	if err != nil {
		panic(err)
	}
	env.Api.RefreshTokenPublicKey = string(bytes)

	return &env
}
