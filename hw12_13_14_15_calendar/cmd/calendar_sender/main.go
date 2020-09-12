package main

import (
	"flag"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/sender"
	"github.com/rs/zerolog/log"
)

var cfgPath string

func fatal(err error) {
	log.Fatal().Err(err).Msg("Application cannot start")
}

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar Sender config")
}

func main() {
	flag.Parse()

	cfg, err := config.NewSender(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)
	if err = logger.Init(&cfg.LogConfig); err != nil {
		fatal(err)
	}

	qCfg := cfg.QueueConfig
	consumer := queue.NewConsumer(
		qCfg.URI,
		qCfg.QueueName,
		qCfg.ExchangeType,
		qCfg.QosPrefetchCount,
		qCfg.MaxReconnectAttempts,
		qCfg.ReconnectTimeoutMs,
	)
	app := sender.New(consumer, qCfg.ScanTimeoutMs)
	if err := app.Run(); err != nil {
		fatal(err)
	}
}
