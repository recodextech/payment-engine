package avro

import (
	"fmt"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	"github.com/recodextech/container"

	"github.com/tryfix/log"
	"github.com/tryfix/schemaregistry"
)

const versionAll = -2

type Encoder interface {
	Subject() string
	JSONDecoder() func(in []byte) (interface{}, error)
}

type SchemaRegistry struct {
	registry   *schemaregistry.Registry
	config     *SchemaRegistryConfig
	baseLogger log.Logger
}

func (s *SchemaRegistry) Init(container container.Container) error {
	s.baseLogger = container.Resolve(application.ModuleBaseLogger).(log.Logger).
		NewLog(log.Prefixed("avro.schema-registry"))
	s.config = container.GetGlobalConfig(application.ModuleSchemaRegistry).(*SchemaRegistryConfig)

	// initiate and assign schema registry
	r, err := schemaregistry.NewRegistry(
		s.config.URL,
		schemaregistry.WithBackgroundSync(s.config.SyncBrokers, s.config.SyncTopic),
		schemaregistry.WithLogger(s.baseLogger))
	if err != nil {
		return errors.Wrap(err, "schema registry init failed")
	}
	s.registry = r

	// register avro encoder in schema registry and bind them into the container
	err = s.registerAndBind(s.encoders(), container)
	if err != nil {
		return errors.Wrap(err, "schema registry init failed")
	}

	// start schema registry background sync
	if err := r.Sync(); err != nil {
		return errors.Wrap(err, "schema registry background sync failed")
	}

	return nil
}

func (s *SchemaRegistry) encoders() map[string]Encoder {
	return map[string]Encoder{
		// schemas
	}
}

func (s *SchemaRegistry) registerAndBind(encoders map[string]Encoder, container container.Container) error {
	for modName, encoder := range encoders {
		sub := encoder.Subject()
		version := s.config.Subjects[encoder.Subject()]
		err := s.registry.Register(sub, versionAll, encoder.JSONDecoder())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(`schema-registry: schema %s(%d) register failed`, sub, version))
		}

		container.Bind(modName, s.registry.WithSchema(encoder.Subject(), s.config.Subjects[encoder.Subject()]))
	}

	return nil
}
