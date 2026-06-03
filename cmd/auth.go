package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required min=8"`
	Role string `json:"role" binding:"required"`
}



func(a *application) register(c *gin.Context) {
	context := c.Request.Context()
	var user RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emailClean := strings.ToLower(strings.TrimSpace(user.Email))
	existingUser, err := a.models.Users.GetByEmail(context, emailClean)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform database query"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	roleClean := strings.ToLower(strings.TrimSpace(user.Role))
	if roleClean != "reporter" && roleClean != "supervisor" && roleClean != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	newUser, err := a.models.Users.Insert(context, user.Name, emailClean, hashedPassword, roleClean)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
	
}