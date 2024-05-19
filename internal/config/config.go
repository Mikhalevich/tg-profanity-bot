package config

type Config struct {
	LogLevel  string    `yaml:"log_level" required:"true"`
	Bot       Bot       `yaml:"bot" required:"true"`
	Profanity Profanity `yaml:"profanity" required:"true"`
}

type Consumer struct {
	LogLevel  string    `yaml:"log_level" required:"true"`
	BotToken  string    `yaml:"bot_token" required:"true"`
	Profanity Profanity `yaml:"profanity" required:"true"`
	Rabbit    RabbitMQ  `yaml:"rabbit" required:"true"`
}

type Bot struct {
	Token                string `yaml:"token" required:"true"`
	UpdateTimeoutSeconds int    `yaml:"update_timeout_seconds" default:"5"`
}

type Profanity struct {
	Dynamic string `yaml:"dynamic"`
	Static  string `yaml:"static"`
}

type RabbitMQ struct {
	URL          string `yaml:"url" required:"true"`
	MsgQueue     string `yaml:"msg_queue" required:"true"`
	WorkersCount int    `yaml:"workers_count" required:"true"`
}
