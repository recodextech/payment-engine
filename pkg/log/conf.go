package log

import (
	"payment-engine/internal/domain/adaptors"
	"payment-engine/pkg/env"
	"payment-engine/pkg/errors"
	defaultLog "log"
)

type LoggerConf struct {
	Colors   bool              `env:"LOG_COLORS" envDefault:"true"`
	FilePath bool              `env:"LOG_FILEPATH_ENABLED" envDefault:"true"`
	Prefix   string            `env:"LOG_PREFIX" envDefault:"application"`
	Level    adaptors.LogLevel `env:"LOG_LEVEL" envDefault:"TRACE"`
}

func (l *LoggerConf) Register() error {
	err := env.Parse(l)
	if err != nil {
		return errors.Wrap(err, "logger: error parsing logger config")
	}

	return nil
}

func (l *LoggerConf) Validate() error {
	return nil
}

func (l *LoggerConf) Print() interface{} {
	defer defaultLog.Println("logger configs loaded")
	return *l
}
