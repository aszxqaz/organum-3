package main

import (
	"organum/internal/jsonlog"
	"organum/internal/store"
	"os"

	"github.com/olahol/melody"
)

type config struct {
	port   string
	secret string
}

type application struct {
	config *config
	logger *jsonlog.Logger
	store  *store.Store
	melody *melody.Melody
}

func main() {
	config := getConfig()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	store := store.NewStore()

	app := &application{
		config,
		logger,
		store,
		melody.New(),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func getConfig() *config {
	port := getEnv("PORT", "5000")
	secret := getEnv("SECRET", "secret")
	return &config{
		port:   port,
		secret: secret,
	}
}

func getEnv(name string, fallback string) string {
	val := os.Getenv(name)
	if val == "" {
		return fallback
	}
	return val
}
