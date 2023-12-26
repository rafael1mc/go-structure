package env

type Environment struct {
	Name    string `envconfig:"ENV_NAME" default:"LOCAL"`
	IsDebug bool   `envconfig:"APP_DEBUG" default:"false"`
}
