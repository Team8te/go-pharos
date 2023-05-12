package main

import (
	"github.com/go-pharos/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	log.SetFormatter(&log.JSONFormatter{})

	app, err := app.NewApp()
	if err != nil {
		return
	}

	err = app.Run()
}
