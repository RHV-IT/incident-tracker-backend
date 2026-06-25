package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentManagementModel struct {
	DB *pgxpool.Pool
}

type IncidentManagement struct {
	Id int `json:"id"`
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
	query := `
		INSERT INTO incident_management (
			incident_id, impact_on_service, contributory_factors, actions_taken_outcomes, recommendations, lessons_learned,
			informed_patient, informed_relative, informed_senior_manager, informed_pharmacist, police_incident_number, informed_other,
			risk_severity, risk_likelihood, risk_rating,
			ohs_absence_over_3_days, ohs_act_of_violence_or_danger, ohs_hospitalisation_over_24_hours, ohs_staff_name, ohs_staff_dob, ohs_staff_address,
			manager_name, manager_signature, manager_designation, manager_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
		RETURNING id;
	`
	err := m.DB.QueryRow(ctx, query, incident.IncidentId, incident.ImpactOnService, incident.ContributoryFactors, incident.ActionsTakenOutcomes, incident.Recommendations, incident.LessonsLearned,
		incident.InformedPatient, incident.InformedRelative, incident.InformedSeniorManager, incident.InformedPharmacist, incident.PoliceIncidentNumber, incident.InformedOther,
		incident.RiskSeverity, incident.RiskLikelihood, incident.RiskRating,
		incident.OhsAbsenceOver3Days, incident.OhsActOfViolenceOrDanger, incident.OhsHospitalisationOver24Hours, incident.OhsStaffName, incident.OhsStaffDob, incident.OhsStaffAddress,
		incident.ManagerName, incident.ManagerSignature, incident.ManagerDesignation, incident.ManagerDate).Scan(&incident.Id)
	if err != nil {
		return IncidentManagement{}, fmt.Errorf("database query error: %w", err)
	}
	return *incident, nil
}