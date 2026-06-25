package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) submitIncidentManagement(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be an admin"})
		return
	}
}