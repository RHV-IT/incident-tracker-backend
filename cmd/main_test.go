package main

import (
	"context"
	"fmt"
	"issueTracking/internal/db"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func mockAuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userRole", role)
		c.Next()
	}
}

func mockEmailMiddleware(email string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userEmail", email)
		c.Next()
	}
}

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

func insertIncident(payload *db.Incident, a *application, t *testing.T) error {
	if _, err := a.models.Incidents.Insert(context.Background(), payload); err != nil {
		return fmt.Errorf("error seeding test data into incidents")
	}

	return nil
}

func insertUser(a *application, t *testing.T) error {
	if _, err := a.models.Users.Insert(context.Background(), "testuser", "testuser@example.com", "testpassword", "admin", "it"); err != nil {
		return fmt.Errorf("error seeding testdata into users: %v", err)
	}
	return nil
}
