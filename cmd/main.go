package main

import (
	"issueTracking/internal/db"
	"issueTracking/internal/env"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)



type application struct {
	port int
	jwtsecret string
	db *pgxpool.Pool
	models db.Models
}

func main() {
	pool, err := db.InitPool()
	if err != nil {
		log.Fatalf("Failed to initialize database connection pool: %v", err)
	}
	defer pool.Close()
	models := db.NewModels(pool)
	app := &application{
		port: env.GetEnvInt("PORT", 3001),
		jwtsecret: env.GetEnvString("jwtSecret", "someSecret"),
		db: pool,
		models: models,
	}
	app.serve()
	
}