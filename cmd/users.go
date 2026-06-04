package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func(a *application) update(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Must be a superadmin"})
		return
	}
}

func(a *application) disable(c *gin.Context) {

}