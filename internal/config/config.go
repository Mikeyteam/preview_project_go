package config

import (
	"github.com/Mikeyteam/preview_project_go/internal/helpers"
	"github.com/caarlos0/env"
)

type Config struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"debug"`
	HTTPListen       string `env:"HTTP_LISTEN" envDefault:":8013"`
	ImageMaxFileSize int    `env:"IMAGE_MAX_FILE_SIZE" envDefault:"1000000"`
	ImageGetTimeout  int    `env:"IMAGE_GET_TIMEOUT" envDefault:"10"`
	CacheSize        int    `env:"CACHE_SIZE" envDefault:"100000000"`
	CacheType        string `env:"CACHE_TYPE" envDefault:"disk"`
	CachePath        string `env:"CACHE_PATH" envDefault:"./cache"`
}

func ConfigFromEnv() Config {
	c := Config{}

	err := env.Parse(&c)
	helpers.FailOnError(err, "fail get config from Env")

	return c
}
