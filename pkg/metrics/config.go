package metrics

import (
	"payment-engine/pkg/env"
	"payment-engine/pkg/errors"
	"log"
	"strings"
	"time"
)

type Config struct {
	System    string `env:"METRICS_SYSTEM" envDefault:"hailing"`
	SubSystem string `env:"METRICS_SUBSYSTEM" envDefault:"golang_skeleton"`
	HTTP      struct {
		Path         string        `env:"METRICS_HTTP_PATH" envDefault:"metrics"`
		Host         string        `env:"METRICS_HTTP_HOST" envDefault:":7001"`
		ShutdownWait time.Duration `env:"METRICS_HTTP_SERVER_SHUTDOWN_WAIT" envDefault:"5s"`
	}
	Timeouts struct {
		Read         time.Duration `env:"HTTP_SERVER_TIMEOUTS_READ" envDefault:"10s"`
		ReadHeader   time.Duration `env:"HTTP_SERVER_TIMEOUTS_READ_HEADER" envDefault:"10s"`
		Write        time.Duration `env:"HTTP_SERVER_TIMEOUTS_WRITE" envDefault:"10s"`
		Idle         time.Duration `env:"HTTP_SERVER_TIMEOUTS_IDLE" envDefault:"10s"`
		ShutdownWait time.Duration `env:"HTTP_SERVER_TIMEOUTS_SHUTDOWN" envDefault:"5s"`
	}
}

func (r *Config) Register() error {
	err := env.Parse(r)
	if err != nil {
		return errors.Wrap(err, "metrics: error loading metrics config")
	}
	return nil
}

func (r *Config) Validate() error {
	if strings.ContainsAny(r.System+r.SubSystem, "-.") {
		return errors.New("METRICS_SYSTEM and METRICS_SUBSYSTEM variables cannot contain special characters," +
			"_,.")
	}
	return nil
}

func (r *Config) Print() interface{} {
	defer log.Println("application configs loaded")
	return *r
}
