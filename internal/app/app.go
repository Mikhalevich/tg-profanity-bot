package app

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/configor"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/storage/postgres"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

func MakeProfanityReplacer(cfg config.Profanity, m profanity.Matcher) processor.TextReplacer {
	if cfg.Dynamic != "" {
		return profanity.New(m, replacer.NewDynamic(cfg.Dynamic))
	}

	return profanity.New(m, replacer.NewStatic(cfg.Static))
}

func MakeMatcher(chatWordsProvider matcher.ChatWordsProvider) (profanity.Matcher, error) {
	words, err := config.BadWords()
	if err != nil {
		return nil, fmt.Errorf("get bad words: %w", err)
	}

	if chatWordsProvider != nil {
		return matcher.NewNewAhocorasickDynamic(chatWordsProvider, words), nil
	}

	return matcher.NewAhocorasick(words), nil
}

func MakeMsgProcessor(
	profanityCfg config.Profanity,
	botToken string,
	chatWordsProvider matcher.ChatWordsProvider,
) (bot.MessageProcessor, error) {
	m, err := MakeMatcher(chatWordsProvider)
	if err != nil {
		return nil, fmt.Errorf("make matcher: %w", err)
	}

	replacer := MakeProfanityReplacer(profanityCfg, m)

	msgSender, err := msgsender.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("make msg sender: %w", err)
	}

	return processor.New(replacer, msgSender), nil
}

func InitPostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	if cfg.Connection == "" {
		return nil, func() {}, nil
	}

	db, err := otelsql.Open("pgx", cfg.Connection)
	if err != nil {
		return nil, func() {}, fmt.Errorf("open database: %w", err)
	}

	p := postgres.New(db, "pgx")

	return p, func() {
		db.Close()
	}, nil
}

// MakeRabbitAMQPChannel make rabbitmq channel and returns channel itself, clearing func and error.
func MakeRabbitAMQPChannel(url string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("ampq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("get channel from connection: %w", err)
	}

	return ch, func() {
		ch.Close()
		conn.Close()
	}, nil
}

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "consumer worker config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger(lvl string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger, nil
}
