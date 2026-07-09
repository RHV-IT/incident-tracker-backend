package main

import (
	"net/http"
	"strconv"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
)

// addComment posts a comment on an incident. Requires admin or manager role.
	userRole := c.GetString("userRole")
	if userRole != "manager" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed add comments to this incident"})
		return
	}
	var comment *db.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context := c.Request.Context()
	if err := a.models.Comments.InsertComment(context, comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "comment added"})
}

// getComments returns comments for a given incident. Requires admin or manager role.
	userRole := c.GetString("userRole")
	if userRole != "admin" && userRole != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to view incident comments"})
		return
	}

	ctx := c.Request.Context()
	incidentId, err := strconv.Atoi(c.DefaultQuery("incidentId", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id was passed"})
		return
	}

	comments, err := a.models.Comments.GetComments(ctx, incidentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
