package main

import (
	"net/http"
	"strconv"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
)

func (a *application) deathReport(c *gin.Context) {
	var deathReport db.DeathReport
	if err := c.ShouldBindJSON(&deathReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A bad request was sent"})
		return
	}
	context := c.Request.Context()
	err := a.models.DeathReport.InsertDeathReport(context, &deathReport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The death has been reported"})
}

func (a *application) updateDeathReport(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter was passed"})
		return
	}
}
