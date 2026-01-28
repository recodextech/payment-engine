package encoders

import (
	"payment-engine/internal/adaptors/encoders/avro"
	"payment-engine/internal/adaptors/encoders/json"
	"payment-engine/internal/adaptors/encoders/misc"
	"payment-engine/internal/domain/application"

	"github.com/recodextech/container"
)

type Encoders struct{}

func (e *Encoders) Init(c container.Container) error {
	// bind encoder modules
	c.Bind(application.ModuleSchemaRegistry, new(avro.SchemaRegistry))
	c.Bind(application.ModuleStringEncoder, new(misc.StringEncoder))
	c.Bind(application.ModuleUUIDEncoder, new(misc.UUIDEncoder))
	c.Bind(application.ModuleJSONEncoder, new(json.Encoder))

	// init schema registry
	// c.Init(application.ModuleSchemaRegistry)

	return nil
}
