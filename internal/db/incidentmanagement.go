package db

import (
	"context"
	"fmt"
	"time"

	"issueTracking/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentManagementModel struct {
	DB *pgxpool.Pool
}

type IncidentManagement struct {
	Id         int `json:"id"`
	IncidentId int `json:"incidentId"`

	ImpactOnService      string `json:"impactOnService" binding:"required"`
	ContributoryFactors  string `json:"contributoryFactors" binding:"required"`
	ActionsTakenOutcomes string `json:"actionsTakenOutcomes" binding:"required"`
	Recommendations      string `json:"recommendations" binding:"required"`
	LessonsLearned       string `json:"lessonsLearned" binding:"required"`

	InformedPatient       bool   `json:"informedPatient"`
	InformedRelative      bool   `json:"informedRelative"`
	InformedSeniorManager bool   `json:"informedSeniorManager"`
	InformedPharmacist    bool   `json:"informedPharmacist"`
	PoliceIncidentNumber  string `json:"policeIncidentNumber,omitempty"`
	InformedOther         string `json:"informedOther,omitempty"`

	RiskSeverity   int `json:"riskSeverity" binding:"required"`
	RiskLikelihood int `json:"riskLikelihood" binding:"required"`
	RiskRating     int `json:"riskRating" binding:"required"`

	OhsAbsenceOver3Days           bool   `json:"ohsAbsenceOver3Days"`
	OhsActOfViolenceOrDanger      bool   `json:"ohsActOfViolenceOrDanger"`
	OhsHospitalizationOver24Hours bool   `json:"ohsHospitalizationOver24Hours"`
	OhsStaffName                  string `json:"ohsStaffName"`
	OhsStaffDob                   string `json:"ohsStaffDob"`
	OhsStaffAddress               string `json:"ohsStaffAddress"`

	ManagerName        string `json:"managerName" binding:"required"`
	ManagerSignature   bool   `json:"managerSignature" binding:"required"`
	ManagerDesignation string `json:"managerDesignation" binding:"required"`
	ManagerDate        string `json:"managerDate" binding:"required"` // date this was filled
}

type IncidentManagementLogs struct {
	Id         int                `json:"id"`
	IncidentId int                `json:"incidentId"`
	ChangedBy  int                `json:"changedBy"`
	Action     string             `json:"action"`
	OldValue   IncidentManagement `json:"oldValue"`
	NewValue   IncidentManagement `json:"newValue"`
	CreatedAt  time.Time          `json:"createdAt"`
	UserName   string             `json:"userName"`
}

func (m *IncidentManagementModel) SubmitReport(ctx context.Context, incident *IncidentManagement) (IncidentManagement, error) {
	query := `
		INSERT INTO incident_management (
			incident_id, impact_on_service, contributory_factors, actions_taken_outcomes, recommendations, lessons_learned,
			informed_patient, informed_relative, informed_senior_manager, informed_pharmacist, police_incident_number, informed_other,
			risk_severity, risk_likelihood, risk_rating,
			ohs_absence_over_3_days, ohs_act_of_violence_or_danger, ohs_hospitalization_over_24_hours, ohs_staff_name, ohs_staff_dob, ohs_staff_address,
			manager_name, manager_signature, manager_designation, manager_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
		RETURNING id;
	`
	err := m.DB.QueryRow(ctx, query, incident.IncidentId, incident.ImpactOnService, incident.ContributoryFactors, incident.ActionsTakenOutcomes, incident.Recommendations, incident.LessonsLearned,
		incident.InformedPatient, incident.InformedRelative, incident.InformedSeniorManager, incident.InformedPharmacist, incident.PoliceIncidentNumber, incident.InformedOther,
		incident.RiskSeverity, incident.RiskLikelihood, incident.RiskRating,
		incident.OhsAbsenceOver3Days, incident.OhsActOfViolenceOrDanger, incident.OhsHospitalizationOver24Hours, incident.OhsStaffName, incident.OhsStaffDob, incident.OhsStaffAddress,
		incident.ManagerName, incident.ManagerSignature, incident.ManagerDesignation, incident.ManagerDate).Scan(&incident.Id)
	if err != nil {
		return IncidentManagement{}, fmt.Errorf("database query error: %w", err)
	}
	return *incident, nil
}

func (m *IncidentManagementModel) FetchById(ctx context.Context, id int) (*IncidentManagement, error) {
	var incidentmanagement IncidentManagement
	query := `
		SELECT 
			id, incident_id, impact_on_service, contributory_factors, actions_taken_outcomes, recommendations, lessons_learned,
			informed_patient, informed_relative, informed_senior_manager, informed_pharmacist, police_incident_number, informed_other,
			risk_severity, risk_likelihood, risk_rating,
			ohs_absence_over_3_days, ohs_act_of_violence_or_danger, ohs_hospitalization_over_24_hours, ohs_staff_name, ohs_staff_dob, ohs_staff_address,
			manager_name, manager_signature, manager_designation, manager_date
		FROM incident_management 
		WHERE incident_id = $1;
	`
	err := m.DB.QueryRow(ctx, query, id).Scan(
		&incidentmanagement.Id,
		&incidentmanagement.IncidentId,
		&incidentmanagement.ImpactOnService,
		&incidentmanagement.ContributoryFactors,
		&incidentmanagement.ActionsTakenOutcomes,
		&incidentmanagement.Recommendations,
		&incidentmanagement.LessonsLearned,
		&incidentmanagement.InformedPatient,
		&incidentmanagement.InformedRelative,
		&incidentmanagement.InformedSeniorManager,
		&incidentmanagement.InformedPharmacist,
		&incidentmanagement.PoliceIncidentNumber,
		&incidentmanagement.InformedOther,
		&incidentmanagement.RiskSeverity,
		&incidentmanagement.RiskLikelihood,
		&incidentmanagement.RiskRating,
		&incidentmanagement.OhsAbsenceOver3Days,
		&incidentmanagement.OhsActOfViolenceOrDanger,
		&incidentmanagement.OhsHospitalizationOver24Hours,
		&incidentmanagement.OhsStaffName,
		&incidentmanagement.OhsStaffDob,
		&incidentmanagement.OhsStaffAddress,
		&incidentmanagement.ManagerName,
		&incidentmanagement.ManagerSignature,
		&incidentmanagement.ManagerDesignation,
		&incidentmanagement.ManagerDate,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return &incidentmanagement, nil
}

func (m *IncidentManagementModel) UpdateIncidentManagement(ctx context.Context, incidentId, userId int, updateIncident *IncidentManagement) error {
	query := `
	UPDATE incident_management
SET  
    impact_on_service = $1, 
    contributory_factors = $2, 
    actions_taken_outcomes = $3, 
    recommendations = $4, 
    lessons_learned = $5,
    informed_patient = $6, 
    informed_relative = $7, 
    informed_senior_manager = $8, 
    informed_pharmacist = $9, 
    police_incident_number = $10, 
    informed_other = $11,
    risk_severity = $12, 
    risk_likelihood = $13, 
    risk_rating = $14,
    ohs_absence_over_3_days = $15, 
    ohs_act_of_violence_or_danger = $16, 
    ohs_hospitalization_over_24_hours = $17, 
    ohs_staff_name = $18, 
    ohs_staff_dob = $19, 
    ohs_staff_address = $20,
    manager_name = $21, 
    manager_signature = $22, 
    manager_designation = $23, 
    manager_date = $24
WHERE incident_id = $25;`

	oldvalue, err := m.FetchById(ctx, incidentId)
	if err != nil {
		return fmt.Errorf("database query err: %w", err)
	}

	_, err = m.DB.Exec(
		ctx, query,
		updateIncident.ImpactOnService,
		updateIncident.ContributoryFactors,
		updateIncident.ActionsTakenOutcomes,
		updateIncident.Recommendations,
		updateIncident.LessonsLearned,
		updateIncident.InformedPatient,
		updateIncident.InformedRelative,
		updateIncident.InformedSeniorManager,
		updateIncident.InformedPharmacist,
		updateIncident.PoliceIncidentNumber,
		updateIncident.InformedOther,
		updateIncident.RiskSeverity,
		updateIncident.RiskLikelihood,
		updateIncident.RiskRating,
		updateIncident.OhsAbsenceOver3Days,
		updateIncident.OhsActOfViolenceOrDanger,
		updateIncident.OhsHospitalizationOver24Hours,
		updateIncident.OhsStaffName,
		updateIncident.OhsStaffDob,
		updateIncident.OhsStaffAddress,
		updateIncident.ManagerName,
		updateIncident.ManagerSignature,
		updateIncident.ManagerDesignation,
		updateIncident.ManagerDate,
		incidentId,
	)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}

	detachedCtx := context.WithoutCancel(ctx)

	go func(bgCtx context.Context) {
		logQuery := `
		INSERT INTO incident_logs
		(incident_id, changed_by, action, old_value, new_value)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id;
		`
		_, logErr := m.DB.Exec(bgCtx, logQuery, incidentId, userId, "updated", oldvalue, updateIncident)
		if logErr != nil {
			fmt.Printf("Asynchronous audit log failed for incident %d: %v\n", incidentId, logErr)
			logger.ErrorFileLogger.Printf("Asynchronous audit log failed for incident %d: %v", incidentId, logErr)
			return
		}

		logger.UpdateIncidentLogger.Printf("Incident %d updated by user %d", incidentId, userId)
	}(detachedCtx)

	return nil
}

func (m *IncidentManagementModel) GetIncidentManagementLogs(ctx context.Context, incidentId int) ([]IncidentManagementLogs, error) {
	query := `
		SELECT incident_logs.*, users.name
		FROM incident_logs INNER JOIN users ON incident_logs.changed_by=users.id
		WHERE incident_logs.incident_id = $1 
		ORDER by incident_logs.incident_id DESC;
	`
	rows, err := m.DB.Query(ctx, query, incidentId)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var IncidentLogs []IncidentManagementLogs
	for rows.Next() {
		var inc IncidentManagementLogs
		err := rows.Scan(&inc.Id, &inc.IncidentId, &inc.ChangedBy, &inc.Action, &inc.OldValue, &inc.NewValue, &inc.CreatedAt, &inc.UserName)
		if err != nil {
			return nil, fmt.Errorf("database query error: %v", err)
		}

		IncidentLogs = append(IncidentLogs, inc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database row iteration error: %v", err)
	}

	return IncidentLogs, nil
}
