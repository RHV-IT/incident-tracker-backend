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

func TestGetIncidentLogsInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPool := db.SetupTestDB(t)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload := map[string]string{
		"test": "test",
	}
	jsonBody, _ := json.Marshal(&payload)

	r := gin.Default()
	r.GET("/api/v1/incidents/:id/managementlogs", mockAuthMiddleware("admin"), a.getIncidentLogs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/incidents/throwerr/managementlogs", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "invalid id parameter was passed", response["error"])
}

func TestGetIncidentLogsInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPool := db.SetupTestDB(t)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	r := gin.Default()
	payload := map[string]string{
		"test": "test",
	}

	jsonBody, _ := json.Marshal(&payload)
	r.GET("/api/v1/incidents/:id/managementlogs", mockAuthMiddleware("manager"), a.getIncidentLogs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/incidents/1/managementlogs", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "you are not allowed to view incident change logs", response["array"])
}
