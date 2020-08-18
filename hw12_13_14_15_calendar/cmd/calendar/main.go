package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/rs/zerolog/log"
)

var ErrUnSupportedRepoType = errors.New("unsupported repository type")

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "", "Calendar config")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Read(cfgPath)
	if err != nil {
		log.Fatal().Err(err)
	}
	err = logger.Init(cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	repo := repository.New(cfg.RepoType)
	if repo == nil {
		log.Fatal().Err(ErrUnSupportedRepoType)
	}
	err = repo.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer repo.Close()

	app, err := calendar.New(repo)
	if err != nil {
		log.Fatal().Err(err)
	}
	err = app.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatal().Err(err)
	}
}
