package db

import "github.com/jackc/pgx/v5/pgxpool"

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"-"`
}

