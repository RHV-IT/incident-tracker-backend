package main

import (
	"net/http"
	"strconv"
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
	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
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
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access users"})
		return
	}
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

func (a *application) getUsers(c *gin.Context) {
	userRole := c.GetString("userRole")
	if userRole != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only super admins can fetch all users"})
		return
	}
	context := c.Request.Context()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit
	users, totalPages, totalItems, err := a.models.Users.GetUsers(context, &limit, &offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, PaginatedUserResponse{
		Data: users,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})
}

func (a *application) userResetPassword(c *gin.Context) {
	context := c.Request.Context()
	userEmail := c.GetString("userEmail")
	var userResetRequest UserResetPassword
	if err := c.ShouldBindJSON(&userResetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	if userResetRequest.Email != userEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to reset another users password"})
		return
	}

	updateHashPassword, err := HashPassword(userResetRequest.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
		return
	}

	if err := a.models.Users.UserResetPassword(context, &userResetRequest.Email, &updateHashPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (a *application) searchUsers(c *gin.Context) {
	searchQuery := c.Query("searchQuery")
	if searchQuery == "" {
		a.getUsers(c)
		return
	}
}
