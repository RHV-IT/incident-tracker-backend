package main

import (
	"issueTracking/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) deathReport(c *gin.Context) {
	var deathReport db.DeathReport
	if err := c.ShouldBindJSON(&deathReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A bad request was sent"})
		return
	}
}
