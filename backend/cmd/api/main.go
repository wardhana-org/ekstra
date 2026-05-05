package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wardhana-org/ekstra/backend/internal/config"
	"github.com/wardhana-org/ekstra/backend/internal/database"
	"github.com/wardhana-org/ekstra/backend/internal/handlers"
	"github.com/wardhana-org/ekstra/backend/internal/repository"
	"github.com/wardhana-org/ekstra/backend/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	// config
	authCookieConfig := handlers.DefaultAuthCookieConfig()
	authCookieConfig.Secure = cfg.AppEnv == "production"

	//router init
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// create repositories
	userRepository := repository.NewUserRepository(db)
	authRepository := repository.NewAuthRepository(db)

	// create services
	authService := services.NewAuthService(userRepository, authRepository)

	// create handlers
	webAuthHandler := handlers.NewWebAuthHandler(authService, authCookieConfig)

	// register routes
	router.GET("/live", handlers.Live())
	router.GET("/ready", handlers.Ready(db))
	router.POST("/auth/login", webAuthHandler.Login)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
