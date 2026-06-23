package main

import (
	"issueTracking/internal/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *application) reportIncident(c *gin.Context) {
	context := c.Request.Context()
	var input db.IncidentReport
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !input.SeverityLevel.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid severity level provided"})
		return
	}

	if !input.IncidentStatus.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid incident status provided"})
		return
	}

	dbIncident := &db.Incident{
		ReporterName:                input.ReporterName,
		Department:                  input.Department,
		Position:                    input.Position,
		ContactInfo:                 input.ContactInfo,
		DateOfIncident:              input.DateOfIncident,
		TimeOfIncident:              input.TimeOfIncident,
		LocationOfIncident:          input.LocationOfIncident,
		TypeOfIncident:              input.TypeOfIncident,
		PeopleInvolved:              input.PeopleInvolved,
		DescriptionOfIncident:       input.DescriptionOfIncident,
		ImmediateActionTaken:        input.ImmediateActionTaken,
		InjuryOrDamage:              input.InjuryOrDamage,
		SeverityLevel:               db.SeverityLevel(input.SeverityLevel),
		SupervisorNotified:          input.SupervisorNotified,
		RecommendedPreventiveAction: input.RecommendedPreventiveAction,
		IncidentStatus:              db.IncidentStatus(input.IncidentStatus),
	}

	savedIncident, err := a.models.Incidents.Insert(context, dbIncident)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute database query"})
		return
	}
	c.JSON(http.StatusOK, savedIncident)
}

func (a *application) getIncidents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit
	context := c.Request.Context()
	userRole := c.GetString("userRole")
	if userRole == "supervisor" || userRole == "reporter" {
		userDepartment := c.GetString("userDepartment")
		incidents, totalPages, totalItems, err := a.models.Incidents.FetchBySupervisor(context, limit, offset, userDepartment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute database query"})
			return
		}
		c.JSON(http.StatusOK, PaginatedIncidentResponse{
			Data: incidents,
			Pagination: PaginationMeta{
				CurrentPage: page,
				PageSize:    limit,
				TotalItems:  totalItems,
				TotalPages:  totalPages,
			},
		})
		return
	}
	incidents, totalPages, totalItems, err := a.models.Incidents.FetchIncidents(context, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute database query"})
		return
	}
	c.JSON(http.StatusOK, PaginatedIncidentResponse{
		Data: incidents,
		Pagination: PaginationMeta{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})
}

func (a *application) updateIncidentStatus(c *gin.Context) {
	context := c.Request.Context()
	userRole := strings.ToLower(c.GetString("userRole"))
	if userRole == "reporter" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update an incident"})
		return
	}
	userDepartment := strings.ToLower(c.GetString("userDepartment"))
	var status IncidentStatusUpdate
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter was passed"})
		return
	}
	fetchedIncident, err := a.models.Incidents.FetchById(context, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}
	incidentDept := strings.ToLower(fetchedIncident.Department)
	if userRole == "supervisor" && userDepartment != incidentDept {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this incident"})
		return
	}
	incident, err := a.models.Incidents.UpdateIncidentStatus(context, id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}

	c.JSON(http.StatusOK, incident)
}
