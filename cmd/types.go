package main

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Role string `json:"role" binding:"required"`
}

type loginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserId int
	Role string
	Email string
	jwt.RegisteredClaims
}