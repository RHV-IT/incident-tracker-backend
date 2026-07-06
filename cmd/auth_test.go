package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestRegisterRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

func mockAuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userRole", role)
		c.Next()
	}
}

func TestRegisterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPool := db.SetupTestDB(t)

	app := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r := gin.Default()
	r.POST("/api/v1/register", mockAuthMiddleware("superadmin"), app.register)
	payload := map[string]string{
		"name":       "test user",
		"email":      "testuser@example.com",
		"password":   "supersecurepassword123",
		"role":       "admin",
		"department": "Engineering",
	}
	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, "testuser@example.com", createdUser["email"])
}

func TestLoginRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPool := db.SetupTestDB(t)

	app := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	hashedPassword, _ := HashPassword("supersecurepassword123")
	ctx := context.Background()
	_, err := app.models.Users.Insert(ctx, "testuse", "testuser@example.com", hashedPassword, "admin", "IT")
	assert.NoError(t, err, "Failed to preseed test user")

	r := app.routes()
	payload := map[string]string{
		"email":    "testuser@example.com",
		"password": "supersecurepassword123",
	}

	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loggedInUser map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &loggedInUser)
	assert.NoError(t, err)
	assert.Equal(t, "testuser@example.com", loggedInUser["email"])
}
