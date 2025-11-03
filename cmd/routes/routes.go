package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zjoart/eunoia/internal/agent"
	"github.com/zjoart/eunoia/internal/checkin"
	"github.com/zjoart/eunoia/internal/config"
	"github.com/zjoart/eunoia/internal/conversation"
	"github.com/zjoart/eunoia/internal/conversation/platforms"
	"github.com/zjoart/eunoia/internal/middleware"
	"github.com/zjoart/eunoia/internal/reflection"
	"github.com/zjoart/eunoia/internal/user"
)

func SetUpRoutes(db *sql.DB, cfg *config.Config) http.Handler {

	allowedOrigins := []string{
		"*",
	}

	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware(allowedOrigins))

	geminiService := agent.NewGeminiService(cfg.AI.GeminiAPIKey)

	userRepo := user.NewRepository(db)
	checkInRepo := checkin.NewRepository(db)
	reflectionRepo := reflection.NewRepository(db)
	conversationRepo := conversation.NewRepository(db)

	conversationService := conversation.NewService(conversationRepo, userRepo, checkInRepo, reflectionRepo, geminiService)

	platform := platforms.NewPlatform("telex")

	conversationHandler := conversation.NewHandler(conversationService, platform)

	router.HandleFunc("/a2a/agent/eunoia", conversationHandler.HandleA2AMessage).Methods("POST")
	router.HandleFunc("/agent/health", conversationHandler.HandleHealthCheck).Methods("GET")
	router.PathPrefix("/.well-known/").Handler(http.StripPrefix("/.well-known/", http.FileServer(http.Dir(".well-known"))))

	return router
}
