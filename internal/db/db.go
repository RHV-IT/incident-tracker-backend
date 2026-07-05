package db

import (
	"context"
	"fmt"
	"time"

	"issueTracking/internal/env"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Users              UserModel
	Incidents          IncidentsModel
	IncidentManagement IncidentManagementModel
	Comments           CommentModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:              UserModel{DB: db},
		Incidents:          IncidentsModel{DB: db},
		IncidentManagement: IncidentManagementModel{DB: db},
		Comments:           CommentModel{DB: db},
	}
}

func InitPool() (*pgxpool.Pool, error) {
	connStr := env.GetEnvString("dbConnStr", "postgres://tracker_user:tracker_password@localhost:5432/issuetracker")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return pool, nil
}
