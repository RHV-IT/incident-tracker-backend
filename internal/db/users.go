package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	Role string `json:"role"`
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := m.DB.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return &user, nil
}

func (m *UserModel) Insert(ctx context.Context, name, email, password, role string) (*User, error) {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, name, email, role`
	var created User
	err := m.DB.QueryRow(ctx, query, name, email, password, role).Scan(
		&created.Id,
		&created.Name,
		&created.Email,
		&created.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return &created, nil
}
