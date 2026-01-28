package http

import (
	"payment-engine/pkg/env"
	"log"
	"time"
)

type Conf struct {
	Host     string `env:"HTTP_SERVER_HOST" envDefault:":80"`
	Timeouts struct {
		Read         time.Duration `env:"HTTP_SERVER_TIMEOUTS_READ" envDefault:"10s"`
		ReadHeader   time.Duration `env:"HTTP_SERVER_TIMEOUTS_READ_HEADER" envDefault:"10s"`
		Write        time.Duration `env:"HTTP_SERVER_TIMEOUTS_WRITE" envDefault:"10s"`
		Idle         time.Duration `env:"HTTP_SERVER_TIMEOUTS_IDLE" envDefault:"10s"`
		ShutdownWait time.Duration `env:"HTTP_SERVER_TIMEOUTS_SHUTDOWN" envDefault:"5s"`
	}
}

func (r *Conf) Register() error {
	return env.Parse(r)
}

func (r *Conf) Validate() error {
	return nil
}

func (r *Conf) Print() interface{} {
	defer log.Println("http configs loaded")
	return *r
}
