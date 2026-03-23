package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func dbInstance() *pgxpool.Pool {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("unable to create connection pool:", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	log.Println("database connected")
	return db
}
