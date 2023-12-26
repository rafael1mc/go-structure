package env

type Api struct {
	Port                   int `envconfig:"API_PORT" default:"8080"`
	AccessTokenTTLSeconds  int `envconfig:"ACCESS_TOKEN_TTL_SECONDS" default:"7200"`    // 2 hours
	RefreshTokenTTLSeconds int `envconfig:"ACCESS_TOKEN_TTL_SECONDS" default:"2592000"` // 30 days
	AccessTokenPrivateKey  string
	AccessTokenPublicKey   string
	RefreshTokenPrivateKey string
	RefreshTokenPublicKey  string

	InviteCodeLength     int `envconfig:"INVITE_CODE_LENGTH" default:"50"`
	InviteCodeTTLSeconds int `envconfig:"INVITE_CODE_TTL_SECONDS" default:"86400"` // 1 day
}
