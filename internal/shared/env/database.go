package env

type Database struct {
	Host string `envconfig:"DB_HOST" default:"localhost"`
	Port int    `envconfig:"DB_PORT" default:"5432"`
	Name string `envconfig:"DB_NAME" default:"postgres"`
	User string `envconfig:"DB_USER" default:"postgres"`
	Pass string `envconfig:"DB_PASS" default:"postgres"`
}
