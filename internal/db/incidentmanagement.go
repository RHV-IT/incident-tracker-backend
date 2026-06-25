package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentManagementModel struct {
	DB *pgxpool.Pool
}

type IncidentManagement struct {
	IncidentId int `json:"incidentId" binding:"required"`

	ImpactOnService string `json:"impactOnService" binding:"required"`
	ContributoryFactors string `json:"contributoryFactors" binding:"required"`
	ActionsTakenOutcomes string `json:"actionsTakenOutcomes" binding:"required"`
	Recommendations string `json:"recommendations" binding:"required"`
	LessonsLearned string `json:"lessonsLearned" binding:"required"`

	InformedPatient bool `json:"informedPatient"`
	InformedRelative bool `json:"informedRelative"`
	InformedSeniorManager bool `json:"informedSeniorManager"`
	InformedPharmacist bool `json:"informedPharmacist"`
	PoliceIncidentNumber string `json:"policeIncidentNumber,omitempty"`
	InformedOther string `json:"informedOther,omitempty"`

	RiskSeverity int `json:"riskSeverity" binding:"required"`
	RiskLikelihood int `json:"riskLikelihood" binding:"required"`
	RiskRating int `json:"riskRating" binding:"required"`

	OhsAbsenceOver3Days bool `json:"ohsAbsenceOver3Days"`
	OhsActOfViolenceOrDanger bool `json:"ohsActOfViolenceOrDanger"`
	OhsHospitalisationOver24Hours bool `json:"ohsHospitalisationOver24Hours"`
	OhsStaffName string `json:"ohsStaffName"`
	OhsStaffDob string `json:"ohsStaffDob"`
	OhsStaffAddress string `json:"ohsStaffAddress"`

	ManagerName string `json:"managerName" binding:"required"`
	ManagerSignature bool `json:"managerSignature" binding:"required"`
	ManagerDesignation string `json:"managerDesignation" binding:"required"`
	ManagerDate string `json:"managerDate" binding:"required"` // date this was filled
}

func(m *IncidentManagementModel) SubmitReport(ctx context.Context, incident *IncidentManagement) (IncidentManagement, error) {
	var incidentManagement IncidentManagement

	return incidentManagement, nil
}