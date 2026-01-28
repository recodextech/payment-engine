package app

import (
	"os"
	"strconv"
)

func DebugMode() bool {
	debug, err := strconv.ParseBool(os.Getenv(`APP_DEBUG`))
	if err != nil {
		return false
	}
	return debug
}
