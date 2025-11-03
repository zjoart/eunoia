package main

import (
	"fmt"
	"net/http"

	"github.com/zjoart/eunoia/cmd/routes"
	"github.com/zjoart/eunoia/internal/config"
	"github.com/zjoart/eunoia/internal/database"
	"github.com/zjoart/eunoia/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found", logger.WithError(err))
	}

	cfg := config.LoadConfig()

	db, errDb := database.InitDB(&cfg.DB)
	if errDb != nil {
		logger.Fatal("Failed to initialize database", logger.WithError(errDb))
	}

	defer db.Close()

	router := routes.SetUpRoutes(db, cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Service starting", logger.Fields{
		"port": cfg.Port,
	})

	// RUN
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Server failed", logger.WithError(err))
	}
}
