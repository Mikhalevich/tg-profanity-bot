package app

import (
	"context"
	"flag"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/configor"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/banprocessor"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/commandstorage"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/mangler"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/mangler/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/mangler/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/permissionchecker"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/staticwords"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/storage/postgres"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func MakeProfanityReplacer(cfg config.Profanity, m mangler.Matcher) port.Mangler {
	if cfg.Dynamic != "" {
		return mangler.New(m, replacer.NewDynamic(cfg.Dynamic))
	}

	return mangler.New(m, replacer.NewStatic(cfg.Static))
}

func MakeMatcher(pg *postgres.Postgres, words []string) mangler.Matcher {
	if pg != nil {
		return matcher.NewNewAhocorasickDynamic(pg, words)
	}

	return matcher.NewAhocorasick(words)
}

func MakeMsgProcessor(
	botToken string,
	pgCfg config.Postgres,
	profanityCfg config.Profanity,
	commandStorageCfg config.CommandRedis,
	banCfg config.BanRedis,
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

	api, err := newBotAPI(botToken)
	if err != nil {
		return nil, nil, fmt.Errorf("create bot api: %w", err)
	}

	commandStorage, err := makeCommandStorage(commandStorageCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("command storage: %w", err)
	}

	banProcessor, err := makeBanProcessor(banCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("ban processor: %w", err)
	}

	return processor.New(
		replacer,
		msgsender.New(api),
		makeWordsProviderFromPG(pg, words),
		makeWordsUpdaterFromPG(pg),
		permissionchecker.New(api),
		commandStorage,
		banProcessor,
		nil,
	), cleanup, nil
}

func newBotAPI(token string) (*tgbotapi.BotAPI, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return api, nil
}

func makeWordsProviderFromPG(pg *postgres.Postgres, words []string) port.WordsProvider {
	if pg != nil {
		return staticwords.New(words, pg)
	}

	return staticwords.New(words, nil)
}

func makeWordsUpdaterFromPG(pg *postgres.Postgres) port.WordsUpdater {
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

func makeCommandStorage(cfg config.CommandRedis) (port.CommandStorage, error) {
	if cfg.Addr == "" {
		return commandstorage.NewNope(), nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return commandstorage.NewRedis(rdb, cfg.TTL), nil
}

func makeBanProcessor(cfg config.BanRedis) (port.BanProcessor, error) {
	if cfg.Addr == "" {
		return banprocessor.NewNope(), nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.DB,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("redis instrument tracing: %w", err)
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return banprocessor.NewRedisBanProcessor(rdb, cfg.BanTTL, cfg.ViolationsPerHour), nil
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

func SetupLogger(lvl string) (logger.Logger, error) {
	l, err := logger.NewLogrusWithLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("creating new logger: %w", err)
	}

	logger.SetStdLogger(l)

	return l, nil
}
