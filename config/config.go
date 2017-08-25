package config

import (
	"os"

	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

// AppConfig contains app's configuration
type AppConfig struct {
	Folder string `default:""`
	Path   string `default:"30d"`

	Google struct {
		Key      string
		Secret   string
		Callback string
	}

	Users []string

	Server struct {
		Port   string `default:"8040"`
		Secret string `default:"Not-So-Secret"`
	}
}

//Config contains global app's configuration
var Config AppConfig

//Load method loads and parses config file
func (c AppConfig) Load(url string) {
	err := configor.Load(&Config, url)
	if err != nil {
		log.Error("Log file not found: " + url)
	}

	if Config.Folder == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		Config.Folder = dir
	}

	log.Info("Serving files from " + Config.Folder)
}
