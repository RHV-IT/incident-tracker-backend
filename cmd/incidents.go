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
		PrincipalName:          input.PrincipalName,
		PrincipalGender:        input.PrincipalGender,
		PrincipalDob:           input.PrincipalDob,
		PrincipalType:          input.PrincipalType,
		PatientId:              input.PatientId,
		PatientWardDept:        input.PatientWardDept,
		StaffJobTitle:          input.StaffJobTitle,
		StaffPhone:             input.StaffPhone,
		StaffPlaceOfWork:       input.StaffPlaceOfWork,
		StaffSite:              input.StaffSite,
		PeopleInvolved:         input.PeopleInvolved,
		DateOfIncident:         input.DateOfIncident,
		TimeOfIncident:         input.TimeOfIncident,
		LocationOfIncident:     input.LocationOfIncident,
		IncidentWardDept:       input.IncidentWardDept,
		Witnesses:              input.Witnesses,
		WitnessType:            input.WitnessType,
		WitnessWardDept:        input.WitnessWardDept,
		WitnessJobTitle:        input.WitnessJobTitle,
		WitnessPhone:           input.WitnessPhone,
		IsNearMiss:             input.IsNearMiss,
		CauseGroup:             input.CauseGroup,
		Causes:                 input.Causes,
		PrescribingDoctor:      input.PrescribingDoctor,
		TreatmentReceived:      input.TreatmentReceived,
		EquipmentInvolved:      input.EquipmentInvolved,
		EquipmentModel:         input.EquipmentModel,
		EquipmentSentForRepair: input.EquipmentSentForRepair,
		EquipmentWithdrawn:     input.EquipmentWithdrawn,
		EquipmentRetained:      input.EquipmentRetained,
		EquipmentNumber:        input.EquipmentNumber,
		IsMedicalDevice:        input.IsMedicalDevice,
		ReporterName:           input.ReporterName,
		ReporterDesignation:    input.ReporterDesignation,
		Signature:              input.Signature,
		ReporterInfo:           input.ReporterInfo,
		ReporterDate:           input.ReporterDate,
		SeverityLevel:          db.SeverityLevel(input.SeverityLevel),
		IncidentStatus:         db.IncidentStatus(input.IncidentStatus),
	}

	savedIncident, err := a.models.Incidents.Insert(context, dbIncident)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
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
	if userRole == "reporter" || userRole == "supervisor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update an incident"})
		return
	}
	userDepartment := strings.ToLower(c.GetString("userDepartment"))
	var status db.IncidentStatusUpdate
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to look up incident"})
		return
	}
	if fetchedIncident == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident report not found"})
		return
	}

	// Scoped matching uses incident_ward_dept to correctly align with clinical spaces
	incidentDept := strings.ToLower(fetchedIncident.IncidentWardDept)
	if userRole == "supervisor" && userDepartment != incidentDept {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this incident"})
		return
	}
	incident, err := a.models.Incidents.UpdateIncidentStatus(context, id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, incident)
}
