package application

import (
	"os"
	"strconv"
)

type Configurator interface {
	Register() error
	Validate() error
	Print() interface{}
}

func DebugMode() bool {
	debug, err := strconv.ParseBool(os.Getenv(`APP_DEBUG`))
	if err != nil {
		return true
	}

	return debug
}
