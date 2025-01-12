package config

import (
	"github.com/joho/godotenv"
	"github.com/maksimulitin/lib/logger"
	"log"
	"log/slog"
)

func LoadConfigEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Error("not found .env file", slog.Any("err", err))
		log.Fatal("Error loading .env file")
	}
}
