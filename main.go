package main

import (
	"io"
	"log"
	"os"

	"github.com/markus621/benchmark/src/app"
	"github.com/markus621/benchmark/src/managers/config"
)

func main() {

	conf := config.New("config.yaml")

	f, _ := os.OpenFile(conf.Logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	log.SetOutput(io.MultiWriter(os.Stderr, f))

	app.Run(conf)
}
