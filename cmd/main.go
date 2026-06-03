package main

import (
	"issueTracking/internal/db"
	"issueTracking/internal/env"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)



type application struct {
	port int
	db *pgxpool.Pool
}

func main() {
	pool, err := db.InitPool()
	if err != nil {
		log.Fatalf("Failed to initialize database connection pool: %v", err)
	}
	app := &application{
		port: env.GetEnvInt("PORT", 3001),
		db: pool,
	}
	defer pool.Close()
	app.serve()
	
}