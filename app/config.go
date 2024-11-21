package app

type Server struct {
	Host string `env:"HOST" env-default:"localhost"`
	Port int    `env:"PORT" env-default:"4000"`
}

type AppConfig struct {
	Server
	AppEnv string `env:"APP_ENV" env-default:"dev"` // "dev", "prodction"
}
