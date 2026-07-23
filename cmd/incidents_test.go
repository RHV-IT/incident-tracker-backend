package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReportIncident(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r := gin.Default()

	r.POST("/api/v1/incidents", a.reportIncident)

	payload := map[string]any{
		"principalName":       "testName",
		"principalGender":     "Male",
		"principalDob":        "today",
		"principalType":       "Patient",
		"patientId":           "iajdaj232",
		"patientWardDept":     "icu",
		"peopleInvolved":      "peopleInvolved",
		"dateOfIncident":      "today",
		"timeOfIncident":      "now",
		"locationOfIncident":  "here",
		"incidentWardDept":    "here?",
		"isNearMiss":          false,
		"causeGroup":          "causeGroup",
		"reporterName":        "Akene Uzezi",
		"reporterDesignation": "???",
		"signature":           true,
		"reporterInfo":        "some info",
		"date":                "today",
		"severityLevel":       "minor",
		"incidentStatus":      "unresolved",
	}
	jsonBody, _ := json.Marshal(&payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/incidents", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "testName", response["principalName"])
}

func TestGetIncidentsDateRange(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	payload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "2026-07-13",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}

	jsonBody, _ := json.Marshal(payload)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertIncident(payload, a, t)
	assert.NoError(t, err)

	r := gin.Default()
	r.GET("/api/v1/incidents", mockAuthMiddleware("admin"), a.getIncidents)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/incidents?dateFrom=2026-07-13&dateTo=2026-07-19", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	t.Logf("response: %+v", response)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetIncidents(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	payload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "today",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}

	jsonBody, _ := json.Marshal(payload)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	err := insertIncident(payload, a, t)
	assert.NoError(t, err)

	r := gin.Default()
	r.GET("/api/v1/incidents", mockAuthMiddleware("admin"), a.getIncidents)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/incidents", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateIncidentStatusForbidden(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r := gin.Default()

	payload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "today",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}

	err := insertIncident(payload, a, t)
	assert.NoError(t, err)

	r.PATCH("/api/v1/incidents/:id/status", mockAuthMiddleware("manager"), a.updateIncidentStatus)

	jsonBody, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/incidents/1/status", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateIncidentStatusInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	r.PATCH("/api/v1/incidents/:id/status", mockAuthMiddleware("admin"), a.updateIncidentStatus)

	payload := &db.IncidentStatusUpdate{
		Status: "resolved",
	}

	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/incidents/badid/status", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateIncidentStatusSuccess(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	r := gin.Default()
	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	payload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "today",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}

	err := insertIncident(payload, a, t)
	assert.NoError(t, err)

	r.PATCH("/api/v1/incidents/:id/status", mockAuthMiddleware("admin"), a.updateIncidentStatus)

	requestPayload := &db.IncidentStatusUpdate{
		Status: "resolved",
	}
	jsonBody, _ := json.Marshal(&requestPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/incidents/1/status", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "resolved", response["incidentStatus"])
}

func TestSubmitIncidentManagement(t *testing.T) {
	db.TruncateTables(t, testPool)

	gin.SetMode(gin.TestMode)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	incidentPayload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "today",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}
	err := insertIncident(incidentPayload, a, t)
	assert.NoError(t, err)

	insertPayload := db.IncidentManagement{
		IncidentId: 3,

		ImpactOnService:      "test",
		ContributoryFactors:  "test",
		ActionsTakenOutcomes: "test",
		Recommendations:      "test",
		LessonsLearned:       "test",

		RiskSeverity:   4,
		RiskLikelihood: 3,
		RiskRating:     3,

		OhsStaffDob:     "test",
		OhsStaffAddress: "test",

		ManagerName:        "test",
		ManagerSignature:   true,
		ManagerDesignation: "test",
		ManagerDate:        "testdate",
	}

	jsonBody, _ := json.Marshal(&insertPayload)

	r := gin.Default()
	r.POST("/api/v1/incidents/:id/management", mockAuthMiddleware("admin"), a.submitIncidentManagement)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/incidents/1/management", bytes.NewBuffer(jsonBody))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchIncidentsNoQuery(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	incidentPayload := &db.Incident{
		PrincipalName:       "testName",
		PrincipalGender:     "Male",
		PrincipalDob:        "today",
		PrincipalType:       "Patient",
		PatientId:           "iajdaj232",
		PatientWardDept:     "icu",
		PeopleInvolved:      "peopleInvolved",
		DateOfIncident:      "today",
		TimeOfIncident:      "now",
		LocationOfIncident:  "here",
		IncidentWardDept:    "here?",
		IsNearMiss:          false,
		CauseGroup:          "causeGroup",
		ReporterName:        "Akene Uzezi",
		ReporterDesignation: "???",
		Signature:           true,
		ReporterInfo:        "some info",
		ReporterDate:        "today",
		SeverityLevel:       "minor",
		IncidentStatus:      "unresolved",
	}

	test, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	err := insertIncident(incidentPayload, a, t)
	assert.NoError(t, err)
	r := gin.Default()
	r.GET("/api/v1/searchIncidents", mockAuthMiddleware("superadmin"), a.searchIncidents)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchIncidents", bytes.NewBuffer(test))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
