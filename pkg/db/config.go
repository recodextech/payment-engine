package db

import (
	"fmt"
	"payment-engine/pkg/env"
	"log"
)

// DatabaseConfig holds db configurations.
type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" envDefault:"localhost"`
	Port     string `env:"DATABASE_PORT" envDefault:"5432"`
	Database string `env:"DATABASE_NAME" envDefault:"public"`
	User     string `env:"DATABASE_USER" envDefault:"dev"`
	Password string `env:"DATABASE_PASSWORD" envDefault:"none"`
}

func (r *DatabaseConfig) Register() error {
	err := env.Parse(r)
	if err != nil {
		return fmt.Errorf("error loading databse Config, %v", err)
	}
	return nil
}

func (r *DatabaseConfig) Validate() error {
	return nil
}

func (r *DatabaseConfig) Print() interface{} {
	defer log.Println("application configs loaded")
	return *r
}
