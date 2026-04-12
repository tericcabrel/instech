package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	toolhttp "tericcabrel/instech/internal/feature/tool/http"
	"tericcabrel/instech/internal/repository"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug") //add -debug flag to enable debug mode

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
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

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Error().Err(pingErr).Msg("Failed to ping database")
		os.Exit(1)
	}

	toolRepository := repository.NewToolRepository(db)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://instech.com", "http://localhost:3000"},
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Instech"))
	})

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Search results for %s", r.URL.Query().Get("q"))))
	})

	toolRouter := toolhttp.InitializeToolRouter(toolRepository)
	r.Mount("/tools", toolRouter)

	log.Info().Msg("Starting server on 8800")

	serverErr := http.ListenAndServe(":8800", r)
	if serverErr != nil {
		log.Error().Err(serverErr).Msg("Failed to start server")
		os.Exit(1)
	}
}
