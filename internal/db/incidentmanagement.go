package db

import "github.com/jackc/pgx/v5/pgxpool"

type IncidentManagementModel struct {
	DB *pgxpool.Pool
}

type IncidentManagement struct {}