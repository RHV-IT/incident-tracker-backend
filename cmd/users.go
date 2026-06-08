package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *application) update(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be a superadmin"})
		return
	}
	context := c.Request.Context()
	var user UpdateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser, err := a.models.Users.GetByEmail(context, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Role = user.Role
	existingUser.Department = user.Department
	updatedUser, err := a.models.Users.Update(context, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (a *application) disable(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be a superadmin"})
		return
	}
	var email DisableRequest
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context := c.Request.Context()
	existingUser, err := a.models.Users.GetByEmail(context, email.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	existingUser.Disabled = true
	updatedUser, err := a.models.Users.DisableUser(context, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (a *application) enable(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized must be a super admin", "role": userRole})
		return
	}
	context := c.Request.Context()
	var email EnableRequest
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser, err := a.models.Users.GetByEmail(context, email.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to perform database query"})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	existingUser.Disabled = false
	updatedUser, err := a.models.Users.EnableUser(context, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to perform database query"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (a *application) getUser(c *gin.Context) {
	context := c.Request.Context()
	userEmail := strings.TrimSpace(c.Query("email"))
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An email address must be passed into the query"})
		return
	}
	user, err := a.models.Users.GetByEmail(context, userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
