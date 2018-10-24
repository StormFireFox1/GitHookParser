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
	log.Info(logrus.Fields{
		"bootTime": time.Now(),
	}, "Booting up...")

	server := server.Server{
		Env: config.New("$HOME/.config/git-hook-parser/config.yml"),
	}

	fmt.Printf("GitHookParser v.%s \n", server.Env.Get("APP_VERSION"))
	fmt.Printf("Listening on port " + fmt.Sprint(server.Env.Get("HANDLER_PORT")))
	server.Routes()

	http.ListenAndServe(fmt.Sprint(server.Env.Get("HANDLER_PORT")), server.Router)
}
