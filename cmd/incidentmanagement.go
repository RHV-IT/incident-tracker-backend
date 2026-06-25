package main

import (
	"issueTracking/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) submitIncidentManagement(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be an admin"})
		return
	}
	var incidentManagement db.IncidentManagement
	if err := c.ShouldBindJSON(&incidentManagement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context := c.Request.Context()
	incidentManagement, err := a.models.IncidentManagement.SubmitReport(context, &incidentManagement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	c.JSON(http.StatusOK, incidentManagement)
}