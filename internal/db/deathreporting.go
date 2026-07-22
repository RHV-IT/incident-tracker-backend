package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeathReportModel struct {
	DB *pgxpool.Pool
}

type DeathReport struct {
	ID                      int    `json:"id"`
	Ref                     string `json:"ref" binding:"required"`
	ReportedDate            string `json:"reportedDate" binding:"required"`
	IncidentDate            string `json:"incidentDate" binding:"required"`
	IncidentTime            string `json:"incidentTime" binding:"required"`
	Department              string `json:"department" binding:"required"`
	Location                string `json:"location" binding:"required"`
	Category                string `json:"category" binding:"required"`
	SubCategory             string `json:"subCategory" binding:"required"`
	Description             string `json:"description" binding:"required"`
	ActionTaken             string `json:"actionTaken" binding:"required"`
	OpenedDate              string `json:"openedDate"`
	SubmittedTime           string `json:"submittedTime"`
	Handler                 string `json:"handler"`
	Manager                 string `json:"manager"`
	Specialty               string `json:"specialty"`
	ExactLocation           string `json:"exactLocation"`
	Coding                  string `json:"coding"`
	Type                    string `json:"type"`
	RiskGrading             string `json:"riskGrading"`
	Result                  string `json:"result"`
	ActualHarm              string `json:"actualHarm"`
	PotentialHarm           string `json:"potentialHarm"`
	Details                 string `json:"details"`
	PatientInvolved         bool   `json:"patientInvolved"`
	PatientTold             bool   `json:"patientTold"`
	FamilyTold              bool   `json:"familyTold"`
	WhatFamilyTold          string `json:"whatFamilyTold"`
	IncidentInvestigation   string `json:"incidentInvestigation"`
	ReviewMeetingDate       string `json:"reviewMeetingDate"`
	QualityAssuranceLead    string `json:"qualityAssuranceLead"`
	DoctorNotified          bool   `json:"doctorNotified"`
	MeetingDiscussionPoints string `json:"meetingDiscussionPoints"`
	MeetingActionPoints     string `json:"meetingActionPoints"`
	LevelOfInvestigation    string `json:"levelOfInvestigation"`
}

func (m *DeathReportModel) InsertDeathReport(ctx context.Context, deathReport *DeathReport) error {
	query := `
    INSERT INTO death_reports (
        ref,
        reported_date,
        incident_date,
        incident_time,
        department,
        location,
        category,
        sub_category,
        description,
        action_taken,
        opened_date,
        submitted_time,
        handler,
        manager,
        specialty,
        exact_location,
        coding,
        type,
        risk_grading,
        result,
        actual_harm,
        potential_harm,
        details,
        patient_involved,
        patient_told,
        family_told,
        what_family_told,
        incident_investigation,
        review_meeting_date,
        quality_assurance_lead,
        doctor_notified,
        meeting_discussion_points,
        meeting_action_points,
        level_of_investigation
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
        $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
        $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
        $31, $32, $33, $34
    )
`
	_, err := m.DB.Exec(
		ctx, query,
		deathReport.Ref,
		deathReport.ReportedDate,
		deathReport.IncidentDate,
		deathReport.IncidentTime,
		deathReport.Department,
		deathReport.Location,
		deathReport.Category,
		deathReport.SubCategory,
		deathReport.Description,
		deathReport.ActionTaken,
		deathReport.OpenedDate,
		deathReport.SubmittedTime,
		deathReport.Handler,
		deathReport.Manager,
		deathReport.Specialty,
		deathReport.ExactLocation,
		deathReport.Coding,
		deathReport.Type,
		deathReport.RiskGrading,
		deathReport.Result,
		deathReport.ActualHarm,
		deathReport.PotentialHarm,
		deathReport.Details,
		deathReport.PatientInvolved,
		deathReport.PatientTold,
		deathReport.FamilyTold,
		deathReport.WhatFamilyTold,
		deathReport.IncidentInvestigation,
		deathReport.ReviewMeetingDate,
		deathReport.QualityAssuranceLead,
		deathReport.DoctorNotified,
		deathReport.MeetingDiscussionPoints,
		deathReport.MeetingActionPoints,
		deathReport.LevelOfInvestigation,
	)
	if err != nil {
		return fmt.Errorf("database query error: %s", err)
	}

	return nil
}
