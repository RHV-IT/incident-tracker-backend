package main

import (
	"bytes"
	"encoding/json"
	"issueTracking/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInsertDeathReportBadRequest(t *testing.T) {
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
