package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config the configuration of this service
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Messaging MessagingConfig `yaml:"messaging"`
}

type ServerConfig struct {
	Host string `yaml:"host" envconfig:"SERVER_HOST"`
	Port string `yaml:"port" envconfig:"SERVER_PORT"`
	Mode string `yaml:"mode" envconfig:"SERVER_MODE"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type" envconfig:"DB_TYPE"`
	Username string `yaml:"user" envconfig:"DB_USERNAME"`
	Password string `yaml:"pass" envconfig:"DB_PASSWORD"`
	Host     string `yaml:"host" envconfig:"DB_HOST"`
	Port     string `yaml:"port" envconfig:"DB_PORT"`
	DbName   string `yaml:"name" envconfig:"DB_NAME"`
}

type MessagingConfig struct {
	Host         string `yaml:"host" envconfig:"MQ_HOST"`
	Port         string `yaml:"port" envconfig:"MQ_PORT"`
	StripTopic   string `yaml:"striptopic" envconfig:"MQ_STRIPTOPIC"`
	ProfileTopic string `yaml:"profiletopic" envconfig:"MQ_STRIPTOPIC"`
	Disabled     bool   `yaml:"disabled" envconfig:"MQ_DISABLED"`
}

// CONFIG the current configuration
// var CONFIG Config

// InitConfig initialize the configuration
func InitConfig(configFile string) Config {
	var cfg Config
	readFile(configFile, &cfg)
	readEnv(&cfg)
	fmt.Printf("%+v", cfg)
	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(configFile string, cfg *Config) {
	abs, err := filepath.Abs(configFile)
	if err != nil {
		processError(err)
	}
	fmt.Printf("Trying to read config from %v", abs)

	f, err := os.Open(configFile)
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
