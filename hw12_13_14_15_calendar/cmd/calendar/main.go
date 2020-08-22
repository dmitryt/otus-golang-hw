package main

import (
	"context"
	"errors"
	"flag"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/rs/zerolog/log"
)

var ErrUnSupportedRepoType = errors.New("unsupported repository type")

var cfgPath string

func fatal(err error) {
	log.Fatal().Err(err).Msg("Application cannot start")
}

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar config")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Read(cfgPath)
	if err != nil {
		fatal(err)
	}
	log.Debug().Msgf("Config inited %+v", cfg)

	if err = logger.Init(cfg); err != nil {
		fatal(err)
	}

	repo := repository.New(cfg.RepoType)
	if repo == nil {
		fatal(ErrUnSupportedRepoType)
	}

	if err = repo.Connect(ctx, cfg); err != nil {
		fatal(err)
	}
	defer repo.Close()

	app, err := app.New(cfg, repo)
	if err != nil {
		fatal(err)
	}

	if err = app.Run(); err != nil {
		fatal(err)
	}
}
