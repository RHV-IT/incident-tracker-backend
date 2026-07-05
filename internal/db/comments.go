package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentModel struct {
	DB *pgxpool.Pool
}

type Comment struct {
	Id         int    `json:"id"`
	IncidentId int    `json:"incidentId" binding:"required"`
	UserId     int    `json:"userId" binding:"required"`
	Comment    string `json:"comment" binding:"required"`
}

func (m *CommentModel) InsertComment(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments
		(incident_id, user_id, comment)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err := m.DB.QueryRow(ctx, query, comment.IncidentId, comment.UserId, comment.Comment).Scan(&comment.Id)
	if err != nil {
		return fmt.Errorf("database execution error: %v", err)
	}

	return nil
}
