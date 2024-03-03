package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config the configuration of this service
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Messaging MessagingConfig `yaml:"messaging"`
	CSV       CSVConfig       `yaml:"csv"`
	Telegram  TelegramConfig  `yaml:"telegram"`
}

type ServerConfig struct {
	Host string `yaml:"host" envconfig:"SERVER_HOST"`
	Port string `yaml:"port" envconfig:"SERVER_PORT"`
	Mode string `yaml:"mode" envconfig:"SERVER_MODE"`
}

type MessagingConfig struct {
	Host         string `yaml:"host" envconfig:"MQ_HOST"`
	Port         string `yaml:"port" envconfig:"MQ_PORT"`
	StripTopic   string `yaml:"striptopic" envconfig:"MQ_STRIPTOPIC"`
	ProfileTopic string `yaml:"profiletopic" envconfig:"MQ_STRIPTOPIC"`
	Disabled     bool   `yaml:"disabled" envconfig:"MQ_DISABLED"`
}
type CSVConfig struct {
	DataDir  string `yaml:"datadir"`
	Interval int    `yaml:"intervalmin"`
}

type TelegramConfig struct {
	Enable         bool    `yaml:"enable" envconfig:"TG_ENABLE"`
	EnableDebug    bool    `yaml:"enabledebug" envconfig:"TG_ENABLE_DEBUG"`
	BotKey         string  `yaml:"apikey" envconfig:"TG_BOT_APIKEY"`
	AllowedUserIDs []int64 `yaml:"allowedusers" envconfig:"TG_ALLOWED_USERS"`
}

// InitConfig initialize the configuration
func InitConfig(configFile string) (cfg *Config, err error) {
	cfg = &Config{}
	abs, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Clean(abs))
	if err != nil {
		return cfg, err
	}
	if err = cfg.readConf(data); err != nil {
		return cfg, err
	}
	log.Printf("%+v", cfg)
	return cfg, nil
}

func (cfg *Config) readConf(data []byte) (err error) {
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return
	}
	return cfg.readEnv()
}

func (cfg *Config) readEnv() (err error) {
	return envconfig.Process("", cfg)
}
