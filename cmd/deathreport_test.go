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

func TestInsertDeathReportBadRequest(t *testing.T) {
	db.TruncateTables(t, testPool)
	payload := map[string]any{
		"reportedDate":            "22/07/2026",
		"incidentDate":            "21/07/2026",
		"incidentTime":            "14:30",
		"department":              "ICU",
		"location":                "Main Hospital",
		"category":                "Mortality",
		"subCategory":             "Unexpected death",
		"description":             "Patient experienced acute cardiac arrest following surgical procedure.",
		"actionTaken":             "CPR initiated immediately; resuscitation team responded. Pronounced dead at 15:05.",
		"openedDate":              "22/07/2026",
		"submittedTime":           "08:00",
		"handler":                 "Dr. John Doe",
		"manager":                 "Jane Smith",
		"specialty":               "Cardiology",
		"exactLocation":           "Bed 4, ICU Ward 2",
		"coding":                  "ICD-10-I46.9",
		"type":                    "Clinical Incident",
		"riskGrading":             "High",
		"result":                  "Fatal",
		"actualHarm":              "Severe / Death",
		"potentialHarm":           "Severe",
		"details":                 "Patient was undergoing routine post-op monitoring.",
		"patientInvolved":         true,
		"patientTold":             false,
		"familyTold":              true,
		"whatFamilyTold":          "Family was informed about cardiac complications and unsuccessful resuscitation efforts.",
		"incidentInvestigation":   "Internal review initiated by QA panel.",
		"reviewMeetingDate":       "25/07/2026",
		"qualityAssuranceLead":    "Dr. Alice Johnson",
		"docNotified":             true,
		"meetingDiscussionPoints": "Reviewed timeline of medication administration and monitoring telemetry logs.",
		"meetingActionPoints":     "Audit telemetry equipment calibration and update post-op cardiac monitoring protocol.",
		"levelOfInvestigation":    "Level 3",
	}
	jsonBody, err := json.Marshal(&payload)
	assert.NoError(t, err)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	r := gin.Default()
	r.POST("/api/v1/deathreport", a.deathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/deathreport", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInsertDeathReportSuccess(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload := map[string]any{
		"ref":                     "DR-2026-001",
		"reportedDate":            "22/07/2026",
		"incidentDate":            "21/07/2026",
		"incidentTime":            "14:30",
		"department":              "ICU",
		"location":                "Main Hospital",
		"category":                "Mortality",
		"subCategory":             "Unexpected death",
		"description":             "Patient experienced acute cardiac arrest following surgical procedure.",
		"actionTaken":             "CPR initiated immediately; resuscitation team responded. Pronounced dead at 15:05.",
		"openedDate":              "22/07/2026",
		"submittedTime":           "08:00",
		"handler":                 "Dr. John Doe",
		"manager":                 "Jane Smith",
		"specialty":               "Cardiology",
		"exactLocation":           "Bed 4, ICU Ward 2",
		"coding":                  "ICD-10-I46.9",
		"type":                    "Clinical Incident",
		"riskGrading":             "High",
		"result":                  "Fatal",
		"actualHarm":              "Severe / Death",
		"potentialHarm":           "Severe",
		"details":                 "Patient was undergoing routine post-op monitoring.",
		"patientInvolved":         true,
		"patientTold":             false,
		"familyTold":              true,
		"whatFamilyTold":          "Family was informed about cardiac complications and unsuccessful resuscitation efforts.",
		"incidentInvestigation":   "Internal review initiated by QA panel.",
		"reviewMeetingDate":       "25/07/2026",
		"qualityAssuranceLead":    "Dr. Alice Johnson",
		"docNotified":             true,
		"meetingDiscussionPoints": "Reviewed timeline of medication administration and monitoring telemetry logs.",
		"meetingActionPoints":     "Audit telemetry equipment calibration and update post-op cardiac monitoring protocol.",
		"levelOfInvestigation":    "Level 3",
	}
	jsonBody, err := json.Marshal(&payload)
	assert.NoError(t, err)

	r := gin.Default()
	r.POST("/api/v1/deathreport", a.deathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/deathreport", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "The death has been reported", response["message"])
}

func TestUpdateRequestInvalidId(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})
	r := gin.Default()
	r.PUT("/deathreport/:id", a.updateDeathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/deathreport/test", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid id parameter was passed", response["error"])
}
