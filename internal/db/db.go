// Package db this is the database package
package db

import (
	"context"
	"fmt"
	"issueTracking/internal/env"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Users              UserModel
	Incidents          IncidentsModel
	IncidentManagement IncidentManagementModel
	Comments           CommentModel
	DeathReport        DeathReportModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:              UserModel{DB: db},
		Incidents:          IncidentsModel{DB: db},
		IncidentManagement: IncidentManagementModel{DB: db},
		Comments:           CommentModel{DB: db},
		DeathReport:        DeathReportModel{DB: db},
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

	var pingErr error
	count := 0
	for i := 0; i < 10; i++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			fmt.Println("database connection established")
			break
		}
		count += 1
		fmt.Println("Failed to connect to Database. Retrying...")
		time.Sleep(800 * time.Millisecond)
		if count == 10 {
			return nil, pingErr
		}
	}

	return pool, nil
}
