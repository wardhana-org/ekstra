package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wardhana-org/ekstra/backend/internal/handlers"
)

func Register(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("/ping", handlers.Ping(db))
}
