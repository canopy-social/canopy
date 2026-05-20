package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Info().
		Str("domain", cfg.Server.Domain).
		Msg("starting canopy worker")

	// Worker implementation will be added in subsequent branches
	// using hibiken/asynq for Redis-backed job processing.
	log.Info().Msg("worker: no handlers registered yet — exiting")
}
