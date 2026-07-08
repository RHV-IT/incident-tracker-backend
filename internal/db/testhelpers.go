//go:test build

package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupTestDBSuite() (*pgxpool.Pool, func(), error) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pgContainer.Terminate(ctx)
	}

	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	var pingErr error

	for i := 0; i < 10; i++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if pingErr != nil {
		cleanup()
		return nil, nil, pingErr
	}

	dir, err := os.Getwd()
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			cleanup()
			return nil, nil, os.ErrNotExist
		}
		dir = parent
	}

	schemaPath := filepath.Join(dir, "tables.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return pool, cleanup, nil
}

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { pgContainer.Terminate(ctx) })

	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	config, err := pgxpool.ParseConfig(connStr)
	assert.NoError(t, err)
	config.MaxConns = 10
	config.MinConns = 2
	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)

	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	require.NoError(t, pingErr, "Postgres failed to become ready in time")

	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root containing go.mod")
		}
		dir = parent
	}

	schemaPath := filepath.Join(dir, "tables.sql")

	schema, err := os.ReadFile(schemaPath)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, string(schema))
	require.NoError(t, err)

	return pool
}

func TruncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE users, incidents, incident_logs, comments RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to truncatetables: %v", err)
	}
}
