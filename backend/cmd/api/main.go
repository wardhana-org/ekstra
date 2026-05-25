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
	authCookieConfig.Domain = cfg.AuthCookieDomain
	authCookieConfig.Secure = cfg.AuthCookieSecure
	authCookieConfig.SameSite = cfg.AuthCookieSameSite

	//router init
	router := gin.Default()
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Fatal("failed to configure trusted proxies: ", err)
	}

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
	router.POST("/auth/register", webAuthHandler.Register)
	router.GET("/auth/me", webAuthHandler.Me)
	router.POST("/auth/refresh", webAuthHandler.Refresh)
	router.POST("/auth/logout", webAuthHandler.Logout)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
