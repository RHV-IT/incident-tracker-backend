package db

import "github.com/jackc/pgx/v5/pgxpool"

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
	DocNotified             bool   `json:"docNotified"`
	MeetingDiscussionPoints string `json:"meetingDiscussionPoints"`
	MeetingActionPoints     string `json:"meetingActionPoints"`
	LevelOfInvestigation    string `json:"levelOfInvestigation"`
}
