package main

import (
	"issueTracking/internal/db"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReportIncident(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r := gin.Default()

	r.POST("/api/v1/incidents", a.reportIncident)
}
