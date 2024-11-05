package main

import (
	"net/url"

	"github.com/Team8te/go-pharos/app"
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

	u11 := url.URL{Scheme: "ws", Host: "127.0.0.1:3000", Path: "/v1/ws"}
	log.Printf("connecting to %s", u11.String())

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	var showHelp bool
	var configPath string
	var client bool
	var uu string
	var path string
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help message")
	pflag.BoolVarP(&client, "client", "", false, "client send")
	pflag.StringVarP(&uu, "url", "", "", "target url")
	pflag.StringVarP(&path, "path", "p", "", "target url")

	pflag.Parse()
	if showHelp {
		pflag.Usage()
		return
	}

	if client {
		var u *url.URL
		u, err = url.Parse(uu)
		if err != nil {
			return
		}

		err = sendContent(u, path)
		return
	}

	app, err := app.NewApp(configPath)
	if err != nil {
		return
	}

	err = app.Run()
}
