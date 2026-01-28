package http

import (
	"payment-engine/pkg/env"
	"log"
)

type KRouterConf struct {
	BootstrapServers []string `env:"HTTP_K_ROUTER_BOOTSTRAP_SERVERS"`
	ApplicationID    string   `env:"HTTP_K_ROUTER_APPLICATION_ID"`
	Topic            string   `env:"HTTP_K_ROUTER_TOPIC"`
	XHeaders         []string `env:"HTTP_EXTERNAL_HEADERS"`
}

func (r *KRouterConf) Register() error {
	err := env.Parse(r)
	if err != nil {
		return err
	}
	if r.Topic == `` {
		r.Topic = `services.hailing.payment-engine.request-rerouted.internal`
	}
	return nil
}

func (r *KRouterConf) Validate() error {
	return nil
}

func (r *KRouterConf) Print() interface{} {
	defer log.Println("krouter configs loaded")
	return *r
}
