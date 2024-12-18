package app

type Server struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port int    `env:"PORT" env-default:"4000"`
}

type AppConfig struct {
	Server
	AppEnv string `env:"APP_ENV" env-default:"dev"` // "dev", "prodction"
}
