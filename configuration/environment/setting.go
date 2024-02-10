package environment

import "time"

type setting struct {
	Application struct {
		ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
	}

	Server struct {
		Context      string        `envconfig:"SERVER_CONTEXT" default:"petshop-system"`
		Port         string        `envconfig:"PORT" default:"9999" required:"true" ignored:"false"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	}

	Postgres struct {
		DBUser     string `envconfig:"DB_USER" default:"petshop-system"`
		DBPassword string `envconfig:"DB_PASSWORD" default:"test1234"`
		DBName     string `envconfig:"DB_NAME" default:"petshop-system"`
		DBHost     string `envconfig:"DB_HOST" default:"localhost"`
		DBPort     string `envconfig:"DB_PORT" default:"5432"`
		DBType     string `envconfig:"DB_TYPE" default:"postgres"`
	}

	RouterConfig struct {
		FileName string `envconfig:"ROUTER_CONFIG" default:"router.json"`
	}
}

var Setting setting
