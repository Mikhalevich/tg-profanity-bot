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
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/staticwords"
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

func MakeMatcher(pg *postgres.Postgres, words []string) profanity.Matcher {
	if pg != nil {
		return matcher.NewNewAhocorasickDynamic(pg, words)
	}

	return matcher.NewAhocorasick(words)
}

func MakeMsgProcessor(
	botToken string,
	pgCfg config.Postgres,
	profanityCfg config.Profanity,
) (bot.MessageProcessor, func(), error) {
	pg, cleanup, err := InitPostgres(pgCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init postgres: %w", err)
	}

	words, err := config.BadWords()
	if err != nil {
		return nil, nil, fmt.Errorf("get bad words: %w", err)
	}

	replacer := MakeProfanityReplacer(profanityCfg, MakeMatcher(pg, words))

	msgSender, err := msgsender.New(botToken)
	if err != nil {
		return nil, nil, fmt.Errorf("make msg sender: %w", err)
	}

	var (
		wordsProvider = makeWordsProviderFromPG(pg, words)
		wordsUpdater  = makeWordsUpdaterFromPG(pg)
	)

	return processor.New(replacer, msgSender, wordsProvider, wordsUpdater), cleanup, nil
}

func makeWordsProviderFromPG(pg *postgres.Postgres, words []string) processor.WordsProvider {
	if pg != nil {
		return pg
	}

	return staticwords.New(words)
}

func makeWordsUpdaterFromPG(pg *postgres.Postgres) processor.WordsUpdater {
	if pg != nil {
		return pg
	}

	return nil
}

func InitPostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	if cfg.Connection == "" {
		return nil, func() {}, nil
	}

	db, err := otelsql.Open("pgx", cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
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
