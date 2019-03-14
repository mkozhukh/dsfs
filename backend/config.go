package main

import (
	"log"
	"os"

	"github.com/jinzhu/configor"
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

	Server struct {
		Port       string `default:"8040"`
		Secret     string `default:"Not-So-Secret"`
		AdminPath  string `default:"/app/admin"`
		ClientPath string `default:"/app/client"`
	}

	DB struct {
		User     string
		Host     string
		Password string
		Database string
	}

	Owner string
}

// Config contains global app's configuration
var Config AppConfig

//Load method loads and parses config file
func (c AppConfig) Load(path string) {
	if path == "" {
		if len(os.Args) > 1 {
			path = os.Args[1]
		} else {
			path = "config.yml"
		}
	}

	err := configor.Load(&Config, path)
	if err != nil {
		log.Printf("Log file not found: %s", path)
	}

	if Config.Folder == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		Config.Folder = dir
	}

	log.Printf("Serving files from " + Config.Folder)
}
