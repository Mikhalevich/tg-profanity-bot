package config

type Config struct {
	BotToken string `yaml:"bot_token" required:"true"`
}
