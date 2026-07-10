package main

import (
	"bytes"
	"encoding/json"
	"issueTracking/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUnauthorized(t *testing.T) {
	db.TruncateTables(t, testPool)
	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r := gin.Default()
	r.PUT("/api/v1/update", mockAuthMiddleware("notsuperadmin"), a.update)

	jsonBody, _ := json.Marshal(&map[string]any{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/update", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Unauthorized. Must be a superadmin", response["error"])
}

func TestUpdateSuccess(t *testing.T) {
	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	err := insertUser(a, t)
	assert.NoError(t, err)

	payload := &UpdateRequest{
		Name:       "testuser",
		Email:      "testuser@example.com",
		Role:       "manager",
		Department: "gopc",
	}
	jsonBody, _ := json.Marshal(&payload)

	r := gin.Default()
	r.PUT("/api/v1/update", mockAuthMiddleware("superadmin"), a.update)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/update", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "manager", response["role"])
}

func TestDisableUnauthorized(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertUser(a, t)
	assert.NoError(t, err)

	payload := &DisableRequest{
		Email: "testuser@example",
	}
	jsonBody, _ := json.Marshal(&payload)

	r := gin.Default()
	r.PUT("/api/v1/disable", mockAuthMiddleware("manager"), a.update)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/disable", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Unauthorized. Must be a superadmin", response["error"])
}
