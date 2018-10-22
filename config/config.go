package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stormfirefox1/GitHookParser/log"
)

// Env is a struct that contains all the information for a server to run by getting the configuration necessary for it.
type Env struct {
	config     *viper.Viper
	configPath string
}

// New returns a config object that can be used by the server struct.
func New(configPath string) Env {
	filePath := strings.Split(configPath, "/")
	configFile := strings.Split(filePath[len(filePath)-1], ".")
	configName := configFile[0]
	configFileType := configFile[1]

	env := Env{
		config:     viper.New(),
		configPath: configPath,
	}

	env.config.SetConfigName(configName)
	env.config.SetConfigType(configFileType)
	env.config.AddConfigPath(strings.Join(filePath[:len(filePath)-1], "/"))

	// setting up default variables for GitHub events (i.e. initial release support)
	githubEvents := make([]string, 4)

	githubEvents[0] = "push"
	githubEvents[3] = "ping"

	// setting sensible defaults for project config
	env.config.SetDefault("HANDLER_PORT", ":80")
	env.config.SetDefault("REDIRECT_URL", "")
	env.config.SetDefault("GITHUB_EVENTS", githubEvents)
	env.config.SetDefault("APP_VERSION", "0.1")

	// allowing for override using envvars
	env.config.AutomaticEnv()

	// reading the config file
	err := env.config.ReadInConfig()

	if err != nil {
		log.Fatal(logrus.Fields{}, fmt.Errorf("Fatal error reading configuration file: %v", err))
	}

	return env
}

// Get will extract a value from the configuration constructs and returns it.
func (e *Env) Get(key string) interface{} {
	return e.config.Get(key)
}
