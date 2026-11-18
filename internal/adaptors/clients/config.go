package clients

import (
	"log"

	"payment-engine/pkg/env"
)

type Config struct {
	NotificationServiceBaseURL string `env:"NOTIFICATION_SERVICE_BASE_URL" envDefault:"http://localhost:8090"`
	HTTPTimeoutInSeconds       int    `env:"HTTP_TIMEOUT_IN_SECONDS" envDefault:"5"`
}

func (r *Config) Register() error {
	return env.Parse(r)
}

func (r *Config) Validate() error {
	return nil
}

func (r *Config) Print() interface{} {
	defer log.Println("http client configs loaded")
	return *r
}
