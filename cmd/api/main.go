package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/zanwyyy/platform/config"
	"github.com/zanwyyy/platform/internal/delivery/http/handler"
	"github.com/zanwyyy/platform/internal/delivery/http/router"
	memrepo "github.com/zanwyyy/platform/internal/repository/memory"
	"github.com/zanwyyy/platform/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Repository layer
	userRepo := memrepo.NewUserRepository()

	// Use case layer
	userUC := usecase.NewUserUseCase(userRepo)

	// Delivery layer
	userHandler := handler.NewUserHandler(userUC)

	engine := gin.Default()
	router.Setup(engine, userHandler)

	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := engine.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
