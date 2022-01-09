package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config the configuration of this service
type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
		Mode string `yaml:"mode" envconfig:"SERVER_MODE"`
	} `yaml:"server"`
	Database struct {
		Type     string `yaml:"type" envconfig:"DB_TYPE"`
		Username string `yaml:"user" envconfig:"DB_USERNAME"`
		Password string `yaml:"pass" envconfig:"DB_PASSWORD"`
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		DbName   string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`
	Messaging struct {
		Host         string `yaml:"host" envconfig:"MQ_HOST"`
		Port         string `yaml:"port" envconfig:"MQ_PORT"`
		StripTopic   string `yaml:"striptopic" envconfig:"MQ_STRIPTOPIC"`
		ProfileTopic string `yaml:"profiletopic" envconfig:"MQ_STRIPTOPIC"`
		Disabled     bool   `yaml:"disabled" envconfig:"MQ_DISABLED"`
	} `yaml:"messaging"`
}

// CONFIG the current configuration
var CONFIG Config

// InitConfig initialize the configuration
func InitConfig() {
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	fmt.Printf("%+v", cfg)
	CONFIG = cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
