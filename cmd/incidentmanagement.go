package main

import (
	"issueTracking/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *application) submitIncidentManagement(c *gin.Context) {
	userRole := c.GetString("userRole")
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter was passed"})
		return
	}
	if userRole != "admin" &&  userRole != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be an admin or supervisor"})
		return
	}
	var incidentManagement db.IncidentManagement
	if err := c.ShouldBindJSON(&incidentManagement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	incidentManagement.IncidentId = id
	context := c.Request.Context()
	incidentManagement, err = a.models.IncidentManagement.SubmitReport(context, &incidentManagement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}
	c.JSON(http.StatusOK, incidentManagement)
}

func (a *application) getIncidentManagement(c *gin.Context) {
	context := c.Request.Context()
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter was passed"})
		return
	}
	incidentManagement, err := a.models.IncidentManagement.FetchById(context, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}

	c.JSON(http.StatusOK, incidentManagement)
}

func (a *application) updateIncidentManagement(c *gin.Context) {
	userRole := c.GetString("userRole")
	context := c.Request.Context()
	if userRole != "supervisor" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized. Must be a supervisor or admin"})
		return
	}
	idParam := c.Param("id")
	incidentId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter was passed"})
		return
	}
	var updateIncident db.IncidentManagement
	uid := c.GetInt("userId")
	if err := c.ShouldBindJSON(&updateIncident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := a.models.IncidentManagement.UpdateIncidentManagement(context, incidentId, uid, &updateIncident); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid id was passed"})
		return
	}
	c.JSON(http.StatusOK, updateIncident)
}
