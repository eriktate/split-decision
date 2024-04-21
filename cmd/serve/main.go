package main

import (
	"fmt"
	"os"

	"github.com/eriktate/splitdecision/bcrypt"
	"github.com/eriktate/splitdecision/http"
	"github.com/eriktate/splitdecision/pg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run() error {
	log.Info().Msg("starting splitdecision...")

	pg, err := pg.New("postgres://sd_appuser:splitd@localhost:5432/splitdecision?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to init connection with postgres: %w", err)
	}

	cfg := http.Config{
		Addr:           "0.0.0.0:9001",
		StaticDir:      "./sd-client/dist/",
		Hasher:         bcrypt.New(10),
		UserService:    pg,
		SessionService: pg,
	}

	return http.Serve(cfg)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if err := run(); err != nil {
		log.Error().Err(err).Msg("splitdecision server failure")
	}
}
