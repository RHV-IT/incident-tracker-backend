package main

import (
	"issueTracking/internal/db"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	pool, cleanup, err := db.SetupTestDBSuite()
	if err != nil {
		log.Fatalf("failed to initialize test container: %v", err)
	}
	testPool = pool

	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}
