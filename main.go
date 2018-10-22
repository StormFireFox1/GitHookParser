package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stormfirefox1/GitHookParser/server"

	"github.com/sirupsen/logrus"

	"github.com/stormfirefox1/GitHookParser/config"
	"github.com/stormfirefox1/GitHookParser/log"
)

func main() {
	fmt.Printf("GitHookParser v.%s \n", config.Get("APP_VERSION"))
	log.Info(logrus.Fields{
		"bootTime": time.Now(),
		"version":  config.Get("APP_VERSION"),
	}, "Booting up...")

	server := server.server{
		env: config.New("$HOME/.config/git-hook-parser"),
	}

	server.routes()

	http.ListenAndServe(config.Get("HANDLER_PORT"), server.router)
}
