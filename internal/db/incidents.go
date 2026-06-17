package db

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentsModel struct {
	DB *pgxpool.Pool
}

type SeverityLevel string

type IncidentStatus string

const (
	NearMiss SeverityLevel = "near miss"
	Minor    SeverityLevel = "minor"
	Major    SeverityLevel = "major"
	Critical SeverityLevel = "critical"
)

const (
	Resolved   IncidentStatus = "resolved"
	InProgress IncidentStatus = "inprogress"
	Unresolved IncidentStatus = "unresolved"
)

func (s SeverityLevel) IsValid() bool {
	switch s {
	case NearMiss, Minor, Major, Critical:
		return true
	}
	return false
}

func (i IncidentStatus) IsValid() bool {
	switch i {
	case Resolved, InProgress, Unresolved:
		return true
	}
	return false
}

type Incident struct {
	Id                          int            `json:"id"`
	ReporterName                string         `json:"reporterName"`
	Department                  string         `json:"department"`
	Position                    string         `json:"position"`
	ContactInfo                 string         `json:"contactInfo"`
	DateOfIncident              string         `json:"dateOfIncident"`
	TimeOfIncident              string         `json:"timeOfIncident"`
	LocationOfIncident          string         `json:"locationOfIncident"`
	TypeOfIncident              string         `json:"typeOfIncident"`
	PeopleInvolved              string         `json:"peopleInvolved"`
	DescriptionOfIncident       string         `json:"descriptionOfIncident"`
	ImmediateActionTaken        string         `json:"immediateActionTaken"`
	InjuryOrDamage              string         `json:"injuryOrDamage"`
	SeverityLevel               SeverityLevel  `json:"severityLevel"`
	SupervisorNotified          string         `json:"supervisorNotified"`
	RecommendedPreventiveAction string         `json:"recommendedPreventiveAction"`
	IncidentStatus              IncidentStatus `json:"incidentStatus"`
}

type IncidentReport struct {
	Id                          int            `json:"id"`
	ReporterName                string         `json:"reporterName"`
	Department                  string         `json:"department"`
	Position                    string         `json:"position"`
	ContactInfo                 string         `json:"contactInfo"`
	DateOfIncident              string         `json:"dateOfIncident"`
	TimeOfIncident              string         `json:"timeOfIncident"`
	LocationOfIncident          string         `json:"locationOfIncident"`
	TypeOfIncident              string         `json:"typeOfIncident"`
	PeopleInvolved              string         `json:"peopleInvolved"`
	DescriptionOfIncident       string         `json:"descriptionOfIncident"`
	ImmediateActionTaken        string         `json:"immediateActionTaken"`
	InjuryOrDamage              string         `json:"injuryOrDamage"`
	SeverityLevel               SeverityLevel  `json:"severityLevel"`
	SupervisorNotified          string         `json:"supervisorNotified"`
	RecommendedPreventiveAction string         `json:"recommendedPreventiveAction"`
	IncidentStatus              IncidentStatus `json:"incidentStatus"`
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
			recommended_preventive_action,
			incident_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, reporter_name, department, position, contact_info, date_of_incident, time_of_incident, location_of_incident, type_of_incident, people_involved, description_of_incident, immediate_action_taken, injury_or_damage, severity_level, supervisor_notified, recommended_preventive_action, incident_status;
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
		incident.IncidentStatus,
	).Scan(
		&incident.Id,
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
		&incident.IncidentStatus,
	)

	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return incident, nil
}

func (m *IncidentsModel) FetchIncidents(ctx context.Context, limit, offset int) ([]IncidentReport, int, int, error) {
	var totalItems int
	err := m.DB.QueryRow(ctx, "SELECT COUNT(*) FROM incidents").Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	query := `
		SELECT
			id,
			reporter_name, department, position, contact_info, 
			date_of_incident, time_of_incident, location_of_incident, 
			type_of_incident, people_involved, description_of_incident, 
			immediate_action_taken, injury_or_damage, severity_level, 
			supervisor_notified, recommended_preventive_action, incident_status 
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
			&inc.Id,
			&inc.ReporterName, &inc.Department, &inc.Position, &inc.ContactInfo,
			&inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident,
			&inc.TypeOfIncident, &inc.PeopleInvolved, &inc.DescriptionOfIncident,
			&inc.ImmediateActionTaken, &inc.InjuryOrDamage, &inc.SeverityLevel,
			&inc.SupervisorNotified, &inc.RecommendedPreventiveAction, &inc.IncidentStatus,
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

func (m *IncidentsModel) FetchBySupervisor(ctx context.Context, limit, offset int, department string) ([]IncidentReport, int, int, error) {
	var totalItems int
	countQuery := `
		SELECT COUNT(*) 
		FROM incidents 
		WHERE LOWER(TRIM(department)) = LOWER(TRIM($1))
	`
	err := m.DB.QueryRow(ctx, countQuery, department).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	query := `
		SELECT 
			id,
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
			recommended_preventive_action,
			incident_status
		FROM incidents 
		WHERE LOWER(TRIM(department)) = LOWER(TRIM($1))
		ORDER BY id DESC 
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
			&inc.Id,
			&inc.ReporterName, &inc.Department, &inc.Position, &inc.ContactInfo,
			&inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident,
			&inc.TypeOfIncident, &inc.PeopleInvolved, &inc.DescriptionOfIncident,
			&inc.ImmediateActionTaken, &inc.InjuryOrDamage, &inc.SeverityLevel,
			&inc.SupervisorNotified, &inc.RecommendedPreventiveAction, &inc.IncidentStatus,
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

func (m *IncidentsModel) FetchById(ctx context.Context, id int) (*IncidentReport, error) {
	query := `
		SELECT
			id,
			reporter_name, department, position, contact_info, 
			date_of_incident, time_of_incident, location_of_incident, 
			type_of_incident, people_involved, description_of_incident, 
			immediate_action_taken, injury_or_damage, severity_level, 
			supervisor_notified, recommended_preventive_action, incident_status 
		FROM incidents
		WHERE id = $1
	`
	var inc IncidentReport
	err := m.DB.QueryRow(ctx, query, id).Scan(
		&inc.Id,
			&inc.ReporterName, &inc.Department, &inc.Position, &inc.ContactInfo,
			&inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident,
			&inc.TypeOfIncident, &inc.PeopleInvolved, &inc.DescriptionOfIncident,
			&inc.ImmediateActionTaken, &inc.InjuryOrDamage, &inc.SeverityLevel,
			&inc.SupervisorNotified, &inc.RecommendedPreventiveAction, &inc.IncidentStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database query error %w", err)
	}

	return &inc, nil
}

func(m *IncidentsModel) UpdateIncidentStatus(context context.Context, id int, status string) (*IncidentReport, error) {
	query := `
		UPDATE incidents
		SET incident_status = $1
		WHERE id = $2`
	if !IncidentStatus(status).IsValid() {
		return nil, fmt.Errorf("invalid incident status")
	}
	_, err := m.DB.Exec(context, query, status, id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	incident, err := m.FetchById(context, id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return incident, nil
}
