package avro

import (
	"encoding/json"
	"payment-engine/pkg/env"
	"payment-engine/pkg/errors"
	defaultLog "log"
)

type Subjects map[string]int

func (s *Subjects) UnmarshalText(bytes []byte) error {
	var subjects map[string]int
	if err := json.Unmarshal(bytes, &subjects); err != nil {
		return err
	}
	*s = subjects

	return nil
}

type SchemaRegistryConfig struct {
	URL         string   `env:"SCHEMA_REGISTRY_URL" envDefault:"localhost:8081"`
	Subjects    Subjects `env:"SCHEMA_REGISTRY_SUBJECTS"`
	SyncEnabled bool     `env:"SCHEMA_REGISTRY_SYNC_ENABLED" envDefault:"false"`
	SyncBrokers []string `env:"SCHEMA_REGISTRY_SYNC_KAFKA_BROKERS" envDefault:"[localhost:9092]"`
	SyncTopic   string   `env:"SCHEMA_REGISTRY_SYNC_SCHEMA_TOPIC" envDefault:"_schemas"`
}

func (s *SchemaRegistryConfig) Register() error {
	err := env.Parse(s)
	if err != nil {
		return errors.Wrap(err, "schema-registry-conf: error parsing schema config")
	}

	return nil
}

func (s *SchemaRegistryConfig) Validate() error {
	return nil
}

func (s *SchemaRegistryConfig) Print() interface{} {
	defer defaultLog.Println("logger configs loaded")

	return *s
}
