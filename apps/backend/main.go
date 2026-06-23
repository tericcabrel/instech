package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"net/http"

	"github.com/joho/godotenv"

	// the _ is used to autoload the sqlite driver
	_ "modernc.org/sqlite"

	"tericcabrel/instech/internal/core"
	"tericcabrel/instech/internal/repository"
)

const (
	SERVER_READ_HEADER_TIMEOUT = 10 * time.Second
	SERVER_READ_TIMEOUT        = 30 * time.Second
	SERVER_WRITE_TIMEOUT       = 30 * time.Second
	SERVER_IDLE_TIMEOUT        = 120 * time.Second
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug") //add -debug flag to enable debug mode

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	//nolint:reassign // zerolog supports this assignment
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("sqlite", databaseURL)
	log.Info().Msgf("Database URL: %s", databaseURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open database")
		os.Exit(1)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Error().Stack().Err(closeErr).Msg("Failed to close database")
		}
	}()

	pingErr := db.PingContext(context.Background())
	if pingErr != nil {
		log.Error().Err(pingErr).Msg("Failed to ping database")
		os.Exit(1)
	}

	toolRepository := repository.NewToolRepository(db)
	relationshipRepository := repository.NewRelationshipRepository(db)

	router := core.HTTPRouter{
		ToolRepository:         toolRepository,
		RelationshipRepository: relationshipRepository,
	}

	handler := router.Initialize()

	log.Info().Msg("Starting server on 8800")

	srv := &http.Server{
		Addr:              ":8801",
		Handler:           handler,
		ReadHeaderTimeout: SERVER_READ_HEADER_TIMEOUT,
		ReadTimeout:       SERVER_READ_TIMEOUT,
		WriteTimeout:      SERVER_WRITE_TIMEOUT,
		IdleTimeout:       SERVER_IDLE_TIMEOUT,
	}
	serverErr := srv.ListenAndServe()
	if serverErr != nil {
		log.Error().Err(serverErr).Msg("Failed to start server")
		os.Exit(1)
	}
}
