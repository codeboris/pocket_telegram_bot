package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Errors    Errors
	Responses Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SaveSuccessfully  string `mapstructure:"save_successfully"`
	UnknowCommand     string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	os.Setenv("TOKEN", "7901022905:AAGVjFPmMdLOQw8NqE9c2cLQyx0YJkp9-d8")
	os.Setenv("CONSUMER_KEY", "112948-6476d81ef67594ff4803305")
	os.Setenv("AUTH_SERVER_URL", "http://localhost/")
	if err := viper.BindEnv("TOKEN"); err != nil {
		return err
	}

	if err := viper.BindEnv("CONSUMER_KEY"); err != nil {
		return err
	}

	if err := viper.BindEnv("AUTH_SERVER_URL"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("TOKEN")
	cfg.PocketConsumerKey = viper.GetString("CONSUMER_KEY")
	cfg.AuthServerURL = viper.GetString("AUTH_SERVER_URL")

	return nil
}
