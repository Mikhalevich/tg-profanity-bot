package config

import (
	"time"
)

type Bot struct {
	LogLevel     string           `yaml:"log_level" required:"true"`
	Tracing      Tracing          `yaml:"tracing" required:"true"`
	Updates      BotUpdates       `yaml:"bot" required:"true"`
	Profanity    Profanity        `yaml:"profanity" required:"true"`
	Rabbit       RabbitMQProducer `yaml:"rabbit"`
	Postgres     Postgres         `yaml:"postgres"`
	CommandRedis CommandRedis     `yaml:"command_redis"`
}

type Consumer struct {
	LogLevel     string           `yaml:"log_level" required:"true"`
	Tracing      Tracing          `yaml:"tracing" required:"true"`
	BotToken     string           `yaml:"bot_token" required:"true"`
	Profanity    Profanity        `yaml:"profanity" required:"true"`
	Rabbit       RabbitMQConsumer `yaml:"rabbit" required:"true"`
	Postgres     Postgres         `yaml:"postgres"`
	CommandRedis CommandRedis     `yaml:"command_redis"`
}

type BotUpdates struct {
	Token                string `yaml:"token" required:"true"`
	UpdateTimeoutSeconds int    `yaml:"update_timeout_seconds" default:"5"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Profanity struct {
	Dynamic string `yaml:"dynamic"`
	Static  string `yaml:"static"`
}

type RabbitMQConsumer struct {
	URL          string `yaml:"url" required:"true"`
	MsgQueue     string `yaml:"msg_queue" required:"true"`
	WorkersCount int    `yaml:"workers_count" required:"true"`
}

type RabbitMQProducer struct {
	URL      string `yaml:"url"`
	MsgQueue string `yaml:"msg_queue"`
}

type Postgres struct {
	Connection string `yaml:"connection"`
}

type CommandRedis struct {
	Addr string        `yaml:"addr"`
	Pwd  string        `yaml:"pwd"`
	DB   int           `yaml:"db"`
	TTL  time.Duration `yaml:"ttl"`
}
