package main

import (
	"github.com/go-pharos/app"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	log.SetFormatter(&log.JSONFormatter{})

	var showHelp bool
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help message")

	pflag.Parse()
	if showHelp {
		pflag.Usage()
		return
	}

	app, err := app.NewApp(configPath)
	if err != nil {
		return
	}

	err = app.Run()
}
