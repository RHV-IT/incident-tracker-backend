package main

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Role string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type UpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Role string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type DisableRequest struct {
	Email string `json:"email" binding:"required"`
}

type EnableRequest struct {
	Email string `json:"email" binding:"required"`
}

type loginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserId int
	Role string
	Email string
	Department string
	jwt.RegisteredClaims
}

type SeverityLevel string

const (
	NearMiss SeverityLevel = "Near Miss"
	Minor SeverityLevel = "Minor"
	Major SeverityLevel = "Major"
	Critical SeverityLevel = "Critical"
)

func (s SeverityLevel) IsValid() bool {
	switch s{
		case NearMiss, Minor, Major, Critical:
			return true
	}
	return  false
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