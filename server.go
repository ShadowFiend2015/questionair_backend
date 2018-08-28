package main

import (
	"flag"
	"questionair_backend/router"
	log "questionair_backend/util/logger"

	l "github.com/labstack/gommon/log"
)

func main() {
	var runmode string
	flag.StringVar(&runmode, "runmode", "prod", "run server at mode")
	flag.Parse()
	log.Logger().SetHeader(`${time_rfc3339}  ${level} ${short_file}  line:${line}`)
	log.Logger().Infof("server run as %s mode", runmode)
	if runmode == "prod" {
		log.Logger().SetLevel(l.WARN)
	}
	router.RunServer(runmode)
}
