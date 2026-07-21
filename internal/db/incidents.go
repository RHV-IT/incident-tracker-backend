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

type (
	SeverityLevel  string
	IncidentStatus string
)

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

// Unified core structure mapping perfectly to the fresh tables.sql schema
type IncidentReport struct {
	Id                     int            `json:"id"`
	PrincipalName          string         `json:"principalName"`
	PrincipalGender        string         `json:"principalGender"`
	PrincipalDob           string         `json:"principalDob"`
	PrincipalType          string         `json:"principalType"` // patient, staff, visiting consultant, other
	PatientId              string         `json:"patientId,omitempty"`
	PatientWardDept        string         `json:"patientWardDept,omitempty"`
	StaffJobTitle          string         `json:"staffJobTitle,omitempty"`
	StaffPhone             string         `json:"staffPhone,omitempty"`
	StaffPlaceOfWork       string         `json:"staffPlaceOfWork,omitempty"`
	StaffSite              string         `json:"staffSite,omitempty"`
	PeopleInvolved         string         `json:"peopleInvolved"`
	DateOfIncident         string         `json:"dateOfIncident"`
	TimeOfIncident         string         `json:"timeOfIncident"`
	LocationOfIncident     string         `json:"locationOfIncident"`
	IncidentWardDept       string         `json:"incidentWardDept"`
	Witnesses              string         `json:"witnesses,omitempty"`
	WitnessType            string         `json:"witnessType,omitempty"`
	WitnessWardDept        string         `json:"witnessWardDept,omitempty"`
	WitnessJobTitle        string         `json:"witnessJobTitle,omitempty"`
	WitnessPhone           string         `json:"witenssPhone,omitempty"` // Preserved frontend JSON tag typo safely
	IsNearMiss             bool           `json:"isNearMiss"`
	CauseGroup             string         `json:"causeGroup"`
	Causes                 string         `json:"causes"`
	PrescribingDoctor      string         `json:"prescribingDoctor"`
	TreatmentReceived      string         `json:"treatmentReceived"`
	EquipmentInvolved      string         `json:"equipmentInvolved"`
	EquipmentModel         string         `json:"equipmentModel,omitempty"`
	EquipmentSentForRepair bool           `json:"equipmentSentForRepair"`
	EquipmentWithdrawn     bool           `json:"equipmentWithdrawn"`
	EquipmentRetained      bool           `json:"equipmentRetained"`
	EquipmentNumber        string         `json:"equipmentNumber,omitempty"`
	IsMedicalDevice        string         `json:"isMedicalDevice,omitempty"`
	ReporterName           string         `json:"reporterName" binding:"required"`
	ReporterDesignation    string         `json:"reporterDesignation" binding:"required"`
	Signature              bool           `json:"signature" binding:"required"`
	ReporterInfo           string         `json:"reporterInfo" binding:"required"`
	ReporterDate           string         `json:"date" binding:"required"`
	SeverityLevel          SeverityLevel  `json:"severityLevel"`
	IncidentStatus         IncidentStatus `json:"incidentStatus"`
}

// Alias Incident to IncidentReport to keep your model layers aligned perfectly
type Incident IncidentReport

type IncidentStatusUpdate struct {
	Status string `json:"incidentStatus" binding:"required"`
}

func (m *IncidentsModel) Insert(ctx context.Context, incident *Incident) (*Incident, error) {
	query := `
		INSERT INTO incidents
		(
			principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39)
		RETURNING id;
	`
	err := m.DB.QueryRow(
		ctx, query,
		incident.PrincipalName, incident.PrincipalGender, incident.PrincipalDob, incident.PrincipalType, incident.PatientId,
		incident.PatientWardDept, incident.StaffJobTitle, incident.StaffPhone, incident.StaffPlaceOfWork, incident.StaffSite,
		incident.PeopleInvolved, incident.DateOfIncident, incident.TimeOfIncident, incident.LocationOfIncident, incident.IncidentWardDept,
		incident.Witnesses, incident.WitnessType, incident.WitnessWardDept, incident.WitnessJobTitle, incident.WitnessPhone,
		incident.IsNearMiss, incident.CauseGroup, incident.Causes, incident.PrescribingDoctor, incident.TreatmentReceived,
		incident.EquipmentInvolved, incident.EquipmentModel, incident.EquipmentSentForRepair, incident.EquipmentWithdrawn, incident.EquipmentRetained,
		incident.EquipmentNumber, incident.IsMedicalDevice, incident.ReporterName, incident.ReporterDesignation, incident.Signature,
		incident.ReporterInfo, incident.ReporterDate, incident.SeverityLevel, incident.IncidentStatus,
	).Scan(&incident.Id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return incident, nil
}

func (m *IncidentsModel) FetchIncidentsByDate(ctx context.Context, limit, offset int, dateFrom, dateTo string) ([]IncidentReport, int, int, error) {
	var totalItems int
	err := m.DB.QueryRow(ctx, "SELECT COUNT(*) FROM incidents WHERE date_of_incident BETWEEN $1 AND $2", dateFrom, dateTo).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}

	query := `
		SELECT
			id, principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
		FROM incidents
		WHERE date_of_incident BETWEEN $1 AND $2
		ORDER BY id DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := m.DB.Query(ctx, query, dateFrom, dateTo, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var incidents []IncidentReport
	for rows.Next() {
		var inc IncidentReport
		err := rows.Scan(
			&inc.Id, &inc.PrincipalName, &inc.PrincipalGender, &inc.PrincipalDob, &inc.PrincipalType, &inc.PatientId,
			&inc.PatientWardDept, &inc.StaffJobTitle, &inc.StaffPhone, &inc.StaffPlaceOfWork, &inc.StaffSite,
			&inc.PeopleInvolved, &inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident, &inc.IncidentWardDept,
			&inc.Witnesses, &inc.WitnessType, &inc.WitnessWardDept, &inc.WitnessJobTitle, &inc.WitnessPhone,
			&inc.IsNearMiss, &inc.CauseGroup, &inc.Causes, &inc.PrescribingDoctor, &inc.TreatmentReceived,
			&inc.EquipmentInvolved, &inc.EquipmentModel, &inc.EquipmentSentForRepair, &inc.EquipmentWithdrawn, &inc.EquipmentRetained,
			&inc.EquipmentNumber, &inc.IsMedicalDevice, &inc.ReporterName, &inc.ReporterDesignation, &inc.Signature,
			&inc.ReporterInfo, &inc.ReporterDate, &inc.SeverityLevel, &inc.IncidentStatus,
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

func (m *IncidentsModel) FetchIncidents(ctx context.Context, limit, offset int) ([]IncidentReport, int, int, error) {
	var totalItems int
	err := m.DB.QueryRow(ctx, "SELECT COUNT(*) FROM incidents").Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}

	query := `
		SELECT
			id, principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
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
			&inc.Id, &inc.PrincipalName, &inc.PrincipalGender, &inc.PrincipalDob, &inc.PrincipalType, &inc.PatientId,
			&inc.PatientWardDept, &inc.StaffJobTitle, &inc.StaffPhone, &inc.StaffPlaceOfWork, &inc.StaffSite,
			&inc.PeopleInvolved, &inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident, &inc.IncidentWardDept,
			&inc.Witnesses, &inc.WitnessType, &inc.WitnessWardDept, &inc.WitnessJobTitle, &inc.WitnessPhone,
			&inc.IsNearMiss, &inc.CauseGroup, &inc.Causes, &inc.PrescribingDoctor, &inc.TreatmentReceived,
			&inc.EquipmentInvolved, &inc.EquipmentModel, &inc.EquipmentSentForRepair, &inc.EquipmentWithdrawn, &inc.EquipmentRetained,
			&inc.EquipmentNumber, &inc.IsMedicalDevice, &inc.ReporterName, &inc.ReporterDesignation, &inc.Signature,
			&inc.ReporterInfo, &inc.ReporterDate, &inc.SeverityLevel, &inc.IncidentStatus,
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

func (m *IncidentsModel) FetchBySupervisorByDate(ctx context.Context, limit, offset int, department, dateFrom, dateTo string) ([]IncidentReport, int, int, error) {
	var totalItems int
	countQuery := `
		SELECT COUNT(*)
		FROM incidents
		WHERE (date_of_incident BETWEEN $2 AND $3)
			AND (LOWER(TRIM(incident_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(patient_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(staff_place_of_work)) = LOWER(TRIM($1))
			)
	`
	err := m.DB.QueryRow(ctx, countQuery, department, dateFrom, dateTo).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}

	query := `
		SELECT
			id, principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
		FROM incidents
		WHERE (date_of_incident BETWEEN $2 AND $3)
			AND (LOWER(TRIM(incident_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(patient_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(staff_place_of_work)) = LOWER(TRIM($1))
			)
		ORDER BY id DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := m.DB.Query(ctx, query, department, dateFrom, dateTo, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var incidents []IncidentReport
	for rows.Next() {
		var inc IncidentReport
		err := rows.Scan(
			&inc.Id, &inc.PrincipalName, &inc.PrincipalGender, &inc.PrincipalDob, &inc.PrincipalType, &inc.PatientId,
			&inc.PatientWardDept, &inc.StaffJobTitle, &inc.StaffPhone, &inc.StaffPlaceOfWork, &inc.StaffSite,
			&inc.PeopleInvolved, &inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident, &inc.IncidentWardDept,
			&inc.Witnesses, &inc.WitnessType, &inc.WitnessWardDept, &inc.WitnessJobTitle, &inc.WitnessPhone,
			&inc.IsNearMiss, &inc.CauseGroup, &inc.Causes, &inc.PrescribingDoctor, &inc.TreatmentReceived,
			&inc.EquipmentInvolved, &inc.EquipmentModel, &inc.EquipmentSentForRepair, &inc.EquipmentWithdrawn, &inc.EquipmentRetained,
			&inc.EquipmentNumber, &inc.IsMedicalDevice, &inc.ReporterName, &inc.ReporterDesignation, &inc.Signature,
			&inc.ReporterInfo, &inc.ReporterDate, &inc.SeverityLevel, &inc.IncidentStatus,
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

	// 1. Update the Count Query to check all three possible department matches
	countQuery := `
		SELECT COUNT(*)
		FROM incidents
		WHERE LOWER(TRIM(incident_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(patient_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(staff_place_of_work)) = LOWER(TRIM($1))
	`
	err := m.DB.QueryRow(ctx, countQuery, department).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("database query error: %w", err)
	}

	// 2. Update the main selection Query with the same OR logic
	query := `
		SELECT
			id, principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
		FROM incidents
		WHERE LOWER(TRIM(incident_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(patient_ward_dept)) = LOWER(TRIM($1))
		   OR LOWER(TRIM(staff_place_of_work)) = LOWER(TRIM($1))
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
			&inc.Id, &inc.PrincipalName, &inc.PrincipalGender, &inc.PrincipalDob, &inc.PrincipalType, &inc.PatientId,
			&inc.PatientWardDept, &inc.StaffJobTitle, &inc.StaffPhone, &inc.StaffPlaceOfWork, &inc.StaffSite,
			&inc.PeopleInvolved, &inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident, &inc.IncidentWardDept,
			&inc.Witnesses, &inc.WitnessType, &inc.WitnessWardDept, &inc.WitnessJobTitle, &inc.WitnessPhone,
			&inc.IsNearMiss, &inc.CauseGroup, &inc.Causes, &inc.PrescribingDoctor, &inc.TreatmentReceived,
			&inc.EquipmentInvolved, &inc.EquipmentModel, &inc.EquipmentSentForRepair, &inc.EquipmentWithdrawn, &inc.EquipmentRetained,
			&inc.EquipmentNumber, &inc.IsMedicalDevice, &inc.ReporterName, &inc.ReporterDesignation, &inc.Signature,
			&inc.ReporterInfo, &inc.ReporterDate, &inc.SeverityLevel, &inc.IncidentStatus,
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
			id, principal_name, principal_gender, principal_dob, principal_type, patient_id,
			patient_ward_dept, staff_job_title, staff_phone, staff_place_of_work, staff_site,
			people_involved, date_of_incident, time_of_incident, location_of_incident, incident_ward_dept,
			witnesses, witness_type, witness_ward_dept, witness_job_title, witness_phone,
			is_near_miss, cause_group, causes, prescribing_doctor, treatment_received,
			equipment_involved, equipment_model, equipment_sent_for_repair, equipment_withdrawn, equipment_retained,
			equipment_number, is_medical_device, reporter_name, reporter_designation, signature,
			reporter_info, reporter_date, severity_level, incident_status
		FROM incidents
		WHERE id = $1
	`
	var inc IncidentReport
	err := m.DB.QueryRow(ctx, query, id).Scan(
		&inc.Id, &inc.PrincipalName, &inc.PrincipalGender, &inc.PrincipalDob, &inc.PrincipalType, &inc.PatientId,
		&inc.PatientWardDept, &inc.StaffJobTitle, &inc.StaffPhone, &inc.StaffPlaceOfWork, &inc.StaffSite,
		&inc.PeopleInvolved, &inc.DateOfIncident, &inc.TimeOfIncident, &inc.LocationOfIncident, &inc.IncidentWardDept,
		&inc.Witnesses, &inc.WitnessType, &inc.WitnessWardDept, &inc.WitnessJobTitle, &inc.WitnessPhone,
		&inc.IsNearMiss, &inc.CauseGroup, &inc.Causes, &inc.PrescribingDoctor, &inc.TreatmentReceived,
		&inc.EquipmentInvolved, &inc.EquipmentModel, &inc.EquipmentSentForRepair, &inc.EquipmentWithdrawn, &inc.EquipmentRetained,
		&inc.EquipmentNumber, &inc.IsMedicalDevice, &inc.ReporterName, &inc.ReporterDesignation, &inc.Signature,
		&inc.ReporterInfo, &inc.ReporterDate, &inc.SeverityLevel, &inc.IncidentStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database query error %w", err)
	}

	return &inc, nil
}

func (m *IncidentsModel) UpdateIncidentStatus(context context.Context, id int, status string) (*IncidentReport, error) {
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
