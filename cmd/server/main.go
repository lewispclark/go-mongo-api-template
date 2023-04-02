package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/go-api-template/pkg/api"
	"github.com/go-api-template/pkg/config"
	"github.com/go-api-template/pkg/engine"
	"github.com/monzo/typhon"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFilename = "config.yaml"
	ListenAddress  = ":8000"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	file, err := os.Open(ConfigFilename)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open config file")
	}

	cfg := &config.Config{}
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config file")
	}

	db, err := cfg.Database.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	e, err := engine.New(&engine.Config{
		Database: db,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create engine")
	}

	rt := api.New(&api.Config{
		Engine: e,
	})

	srv, err := typhon.Listen(rt.Serve(), ListenAddress, typhon.WithTimeout(typhon.TimeoutOptions{Read: time.Second * 10}))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Info().Str("address", ListenAddress).Msg("started listening")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Stop(c)
}
