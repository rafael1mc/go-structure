package env

type Redis struct {
	Host string `envconfig:"REDIS_HOST" default:"localhost"`
	Port int    `envconfig:"REDIS_PORT" default:"6379"`
	Pass string `envconfig:"REDIS_PASS" default:"123456"`
}
