package main

import (
	"issueTracking/internal/db"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type UpdateRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type DisableRequest struct {
	Email string `json:"email" binding:"required"`
}

type EnableRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"` // Added custom password field
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserId     int
	Role       string
	Email      string
	Department string
	jwt.RegisteredClaims
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type PaginatedIncidentResponse struct {
	Data       []db.IncidentReport `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

type IncidentStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
