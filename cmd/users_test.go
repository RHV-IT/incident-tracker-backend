package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issueTracking/internal/db"

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

func TestGetUsersUnauthorized(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertUser(a, t)
	assert.NoError(t, err)

	jsonBody, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/users", mockAuthMiddleware("manager"), a.getUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users?email=testuser@email.com", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUsersSuccess(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertUser(a, t)
	assert.NoError(t, err)

	jsonBody, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/users", mockAuthMiddleware("superadmin"), a.getUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users?email=testuser@email.com", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserResetPassword(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertUser(a, t)
	assert.NoError(t, err)

	jsonBody, _ := json.Marshal(&map[string]any{
		"email":       "testuser@example.com",
		"newPassword": "resetpassword",
	})

	r := gin.Default()
	r.PUT("/api/v1/auth/userResetPassword", mockEmailMiddleware("testuser@example.com"), a.userResetPassword)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/auth/userResetPassword", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsersInvalidRole(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	jsonBody, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/users", mockAuthMiddleware("admin"), a.getUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]any

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "only super admins can fetch all users", response["error"])
}

func TestGetUsers(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	insertUser(a, t)

	jsonBody, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/users", mockAuthMiddleware("superadmin"), a.getUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	var response map[string]any

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	insertUser(a, t)

	payload := &UpdateRequest{
		Name:       "testuser",
		Email:      "testuser@example.com",
		Role:       "superadmin",
		Department: "it",
	}

	jsonBody, _ := json.Marshal(&payload)
	r := gin.Default()
	r.PUT("/api/v1/update", mockAuthMiddleware("superadmin"), a.update)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/update", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	user := response["user"].(map[string]any)

	assert.Equal(t, "superadmin", user["role"])
}

func TestGetUsersNoQuery(t *testing.T) {
	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/searchUsers", mockAuthMiddleware("superadmin"), a.searchUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchUsers", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUsers(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	insertUser(a, t)

	r := gin.Default()
	r.GET("/api/v1/searchUsers", mockAuthMiddleware("superadmin"), a.searchUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchUsers?searchQuery='testuser@example.com'", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
