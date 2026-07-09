package main

import (
	"issueTracking/internal/db"
	"issueTracking/internal/env"
	"issueTracking/internal/logger"
	"log"

	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// application ties together the database connection pool, data models,
// JWT secret, and allowed origins into a single server instance.
type application struct {
	port      int
	jwtsecret string
	db        *pgxpool.Pool
	models    db.Models
	origins   string
}

// main loads the environment, initializes the database pool and logger,
// constructs the application, and starts the HTTP server.
	logger.InitLogger()

	_ = godotenv.Load()

	pool, err := db.InitPool()
	if err != nil {
		log.Fatalf("Failed to initialize database connection pool: %v", err)
	}
	defer pool.Close()
	models := db.NewModels(pool)
	app := &application{
		port:      env.GetEnvInt("PORT", 3001),
		jwtsecret: env.GetEnvString("jwtSecret", "someSecret"),
		db:        pool,
		models:    models,
		origins:   env.GetEnvString("allowedOrigins", "http://localhost:3000,http://192.168.9.227:3000"),
	}
	app.serve()
}
