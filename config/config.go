package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stormfirefox1/GitHookParser/log"
)

var v = viper.New()

func init() {
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("$HOME/.config/git-hook-parser/")

	// setting up default variables for the whole shebang
	githubEvents := make([]string, 4)

	githubEvents[0] = "push"
	githubEvents[1] = "status"
	githubEvents[2] = "pull_request"
	githubEvents[3] = "issues"

	// setting sensible defaults for project config
	v.SetDefault("HANDLER_PORT", ":80")
	v.SetDefault("REDIRECT_URL", "")
	v.SetDefault("GITHUB_EVENTS", githubEvents)

	// allowing for override using envvars
	v.AutomaticEnv()

	// reading the config file
	err := v.ReadInConfig()

	if err != nil {
		log.Fatal(logrus.Fields{}, fmt.Errorf("Fatal error reading configuration file: %v", err))
	}
}

// Get will extract a value from the configuration constructs and returns it.
func Get(key string) interface{} {
	return v.Get(key)
}
