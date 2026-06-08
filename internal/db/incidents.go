package db

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentsModel struct {
	DB *pgxpool.Pool
}

type SeverityLevel string

const (
	NearMiss SeverityLevel = "near miss"
	Minor SeverityLevel = "minor"
	Major SeverityLevel = "major"
	Critical SeverityLevel = "critical"
)

func (s SeverityLevel) IsValid() bool {
	switch s{
		case NearMiss, Minor, Major, Critical:
			return true
	}
	return  false
}

type Incident struct {
	Id int `json:"id"`
	ReporterName string `json:"reporterName"`
	Department string `json:"department"`
	Position string `json:"position"`
	ContactInfo string `json:"contactInfo"`
	DateOfIncident string `json:"dateOfIncident"`
	TimeOfIncident string `json:"timeOfIncident"`
	LocationOfIncident string `json:"locationOfIncident"`
	TypeOfIncident string `json:"typeOfIncident"`
	PeopleInvolved string `json:"peopleInvolved"`
	DescriptionOfIncident string `json:"descriptionOfIncident"`
	ImmediateActionTaken string `json:"immediateActionTaken"`
	InjuryOrDamage string `json:"injuryOrDamage"` 
	SeverityLevel SeverityLevel `json:"severityLevel"`
	SupervisorNotified string `json:"supervisorNotified"`
	RecommendedPreventiveAction string `json:"recommendedPreventiveAction"` 
}

type IncidentReport struct {
	ReporterName string `json:"reporterName"`
	Department string `json:"department"`
	Position string `json:"position"`
	ContactInfo string `json:"contactInfo"`
	DateOfIncident string `json:"dateOfIncident"`
	TimeOfIncident string `json:"timeOfIncident"`
	LocationOfIncident string `json:"locationOfIncident"`
	TypeOfIncident string `json:"typeOfIncident"`
	PeopleInvolved string `json:"peopleInvolved"`
	DescriptionOfIncident string `json:"descriptionOfIncident"`
	ImmediateActionTaken string `json:"immediateActionTaken"`
	InjuryOrDamage string `json:"injuryOrDamage"` 
	SeverityLevel SeverityLevel `json:"severityLevel"`
	SupervisorNotified string `json:"supervisorNotified"`
	RecommendedPreventiveAction string `json:"recommendedPreventiveAction"`
}

func (m *IncidentsModel) Insert(ctx context.Context, incident *Incident) (*Incident, error) {
	query := `
		INSERT INTO incidents
		(
			reporter_name, 
			department, 
			position, 
			contact_info, 
			date_of_incident, 
			time_of_incident, 
			location_of_incident, 
			type_of_incident, 
			people_involved, 
			description_of_incident, 
			immediate_action_taken, 
			injury_or_damage, 
			severity_level, 
			supervisor_notified, 
			recommended_preventive_action
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING reporter_name, department, position, contact_info, date_of_incident, time_of_incident, location_of_incident, type_of_incident, people_involved, description_of_incident, immediate_action_taken, injury_or_damage, severity_level, supervisor_notified, recommended_preventive_action;
	`
	err := m.DB.QueryRow(ctx, query,
		incident.ReporterName,
		incident.Department,
		incident.Position,
		incident.ContactInfo,
		incident.DateOfIncident,
		incident.TimeOfIncident,
		incident.LocationOfIncident,
		incident.TypeOfIncident,
		incident.PeopleInvolved,
		incident.DescriptionOfIncident,
		incident.ImmediateActionTaken,
		incident.InjuryOrDamage,
		incident.SeverityLevel,
		incident.SupervisorNotified,
		incident.RecommendedPreventiveAction,
	).Scan(
		&incident.ReporterName,
		&incident.Department,
		&incident.Position,
		&incident.ContactInfo,
		&incident.DateOfIncident,
		&incident.TimeOfIncident,
		&incident.LocationOfIncident,
		&incident.TypeOfIncident,
		&incident.PeopleInvolved,
		&incident.DescriptionOfIncident,
		&incident.ImmediateActionTaken,
		&incident.InjuryOrDamage,
		&incident.SeverityLevel,
		&incident.SupervisorNotified,
		&incident.RecommendedPreventiveAction,
	)

	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return incident, nil
}

func(m *IncidentsModel) FetchIncidents(ctx context.Context, limit, offset int) ([]IncidentReport, int, int, error) {
	var totalItems int
	err := m.DB.QueryRow(ctx, "SELECT COUNT(*) FROM incidents").Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	query := `
		SELECT 
			reporter_name, department, position, contact_info, 
			date_of_incident, time_of_incident, location_of_incident, 
			type_of_incident, people_involved, description_of_incident, 
			immediate_action_taken, injury_or_damage, severity_level, 
			supervisor_notified, recommended_preventive_action 
		FROM incidents 
		ORDER BY id DESC 
		LIMIT $1 OFFSET $2
	`
	rows, err := m.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()
	var incidents []IncidentReport
	for rows.Next() {
		var inc IncidentReport
		err := rows.Scan(
			&inc.ReporterName, &inc.Department, &inc.Position, &inc.ContactInfo,
			&inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident,
			&inc.TypeOfIncident, &inc.PeopleInvolved, &inc.DescriptionOfIncident,
			&inc.ImmediateActionTaken, &inc.InjuryOrDamage, &inc.SeverityLevel,
			&inc.SupervisorNotified, &inc.RecommendedPreventiveAction,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("database query error: %w", err)
		}
		incidents = append(incidents, inc)
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return incidents, totalPages, totalItems, nil
}

func(m *IncidentsModel) FetchBySupervisor(ctx context.Context, limit, offset int, department string) ([]IncidentReport, int, int, error) {
	var totalItems int
	countQuery := `
SELECT COUNT(*) 
FROM incidents i
JOIN users u ON i.reporterName = u.name 
WHERE u.role = 'supervisor' 
  AND u.department = $1`
	err := m.DB.QueryRow(ctx, countQuery, department).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	query := `
	SELECT 
    i.reporter_name, 
    i.department, 
    i.position, 
    i.contact_info, 
    i.date_of_incident, 
    i.time_of_incident, 
    i.location_of_incident, 
    i.type_of_incident, 
    i.people_involved, 
    i.description_of_incident, 
    i.immediate_action_taken, 
    i.injury_or_damage, 
    i.severity_level, 
    i.supervisor_notified, 
    i.recommended_preventive_action 
FROM incidents i
INNER JOIN users u ON i.reporter_name = u.name 
WHERE u.role = 'supervisor' 
  AND u.department = $1
ORDER BY i.id DESC 
LIMIT $2 OFFSET $3	
	`
	rows, err := m.DB.Query(ctx, query, department, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()
	var incidents []IncidentReport
	for rows.Next() {
		var inc IncidentReport
		err := rows.Scan(
			&inc.ReporterName, &inc.Department, &inc.Position, &inc.ContactInfo,
			&inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident,
			&inc.TypeOfIncident, &inc.PeopleInvolved, &inc.DescriptionOfIncident,
			&inc.ImmediateActionTaken, &inc.InjuryOrDamage, &inc.SeverityLevel,
			&inc.SupervisorNotified, &inc.RecommendedPreventiveAction,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("database query error: %w", err)
		}
		incidents = append(incidents, inc)
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return incidents, totalPages, totalItems, nil
}