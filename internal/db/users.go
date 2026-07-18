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
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Disabled   bool   `json:"disabled"`
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, name, email, password, role, department, disabled FROM users WHERE email = $1`
	err := m.DB.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.Department, &user.Disabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return &user, nil
}

func (m *UserModel) Insert(ctx context.Context, name, email, password, role, department string) (*User, error) {
	query := `INSERT INTO users (name, email, password, role, department) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, email, role, department`
	var created User
	err := m.DB.QueryRow(ctx, query, name, email, password, role, department).Scan(
		&created.Id,
		&created.Name,
		&created.Email,
		&created.Role,
		&created.Department,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return &created, nil
}

func (m *UserModel) Update(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE users SET name = $1, email = $2, role = $3, department = $4 WHERE id = $5`
	_, err := m.DB.Exec(ctx, query, user.Name, user.Email, user.Role, user.Department, user.Id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return user, nil
}

func (m *UserModel) UserResetPassword(ctx context.Context, email, hashedPassword *string) error {
	query := `UPDATE users SET password = $1 WHERE email = $2`
	_, err := m.DB.Exec(ctx, query, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}
	return err
}

func (m *UserModel) DisableUser(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE users SET disabled = $1 WHERE id = $2`
	_, err := m.DB.Exec(ctx, query, true, user.Id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return user, nil
}

func (m *UserModel) EnableUser(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE users SET disabled = $1 WHERE id = $2`
	_, err := m.DB.Exec(ctx, query, false, user.Id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return user, nil
}

func (m *UserModel) ResetPassword(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := m.DB.Exec(ctx, query, user.Password, user.Id)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	return user, nil
}

func (m *UserModel) GetUsers(ctx context.Context) ([]User, error) {
	return nil, nil
}
