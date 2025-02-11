package config

type (
	Config struct {
		HTTP
		Log
		PostgreSQL
		CatAPI
	}

	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}

	PostgreSQL struct {
		User     string `env:"POSTGRESQL_USER"`
		Password string `env:"POSTGRESQL_PASSWORD"`
		Host     string `env:"POSTGRESQL_HOST"`
		Database string `env:"POSTGRESQL_DATABASE"`
	}

	CatAPI struct {
		URL string `env:"CAT_API_URL"`
	}
)
