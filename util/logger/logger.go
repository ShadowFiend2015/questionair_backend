package logger

import (
	"github.com/labstack/gommon/log"
)

var lg *log.Logger

func init() {
	lg = log.New("")
}

func Logger() *log.Logger {
	return lg
}
