package config

type Config struct {
	BotToken             string `yaml:"bot_token" required:"true"`
	UpdateTimeoutSeconds int    `yaml:"update_timeout_seconds" default:"5"`
}
