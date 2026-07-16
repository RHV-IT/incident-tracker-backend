package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (a *application) register(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be a superadmin"})
		return
	}
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

	if user.Password == "" {
		hashedPassword, err := HashPassword("redeemershealthvillage")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		departmentClean := strings.ToLower(strings.TrimSpace(user.Department))

		roleClean := strings.ToLower(strings.TrimSpace(user.Role))
		if roleClean != "reporter" && roleClean != "supervisor" && roleClean != "admin" && roleClean != "superadmin" && roleClean != "manager" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
			return
		}

		newUser, err := a.models.Users.Insert(context, user.Name, emailClean, hashedPassword, roleClean, departmentClean)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add user"})
			return
		}

		c.JSON(http.StatusCreated, newUser)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	departmentClean := strings.ToLower(strings.TrimSpace(user.Department))

	roleClean := strings.ToLower(strings.TrimSpace(user.Role))
	if roleClean != "reporter" && roleClean != "supervisor" && roleClean != "admin" && roleClean != "superadmin" && roleClean != "manager" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	newUser, err := a.models.Users.Insert(context, user.Name, emailClean, hashedPassword, roleClean, departmentClean)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (a *application) login(c *gin.Context) {
	context := c.Request.Context()
	var user loginRequest
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
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if existingUser.Disabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "This account has been disabled"})
		return
	}

	if !CompareHash(user.Password, existingUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}
	claims := &Claims{
		UserId:     existingUser.Id,
		Role:       existingUser.Role,
		Email:      existingUser.Email,
		Department: existingUser.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.jwtsecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": existingUser})
}

func (a *application) resetPassword(c *gin.Context) {
	context := c.Request.Context()
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only super admins are allowed to update a user"})
		return
	}
	var req ResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	emailClean := strings.TrimSpace(req.Email)
	existingUser, err := a.models.Users.GetByEmail(context, emailClean)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing database query"})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	existingUser.Password = hashedPassword

	updatedUser, err := a.models.Users.ResetPassword(context, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}
