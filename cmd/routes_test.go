package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	// create app instance
	a := &application{}

	// get router instance
	router := a.routes()

	// create response recorder
	w := httptest.NewRecorder()

	// create mock http request
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)

	// send requset to router
	router.ServeHTTP(w, req)

	// check if response matches expectation
	assert.Equal(t, http.StatusOK, w.Code)

	// parse and check the json body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "pong", response["message"])
}
