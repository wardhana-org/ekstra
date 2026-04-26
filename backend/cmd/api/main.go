package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wardhana-org/ekstra/backend/internal/config"
	"github.com/wardhana-org/ekstra/backend/internal/database"
	"github.com/wardhana-org/ekstra/backend/internal/routes"
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

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.Register(router, db)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
