package db

import (
	"context"
	"fmt"
	"issueTracking/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentModel struct {
	DB *pgxpool.Pool
}

type Comment struct {
	Id              int    `json:"id"`
	IncidentId      int    `json:"incidentId" binding:"required"`
	UserId          int    `json:"userId" binding:"required"`
	Comment         string `json:"comment" binding:"required"`
	CommentUserName string `json:"commentUserName"`
	CommentUserRole string `json:"commentUserRole"`
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

	go func() {
		logger.CommentLogger.Printf("Comment was made by user: %d", comment.UserId)
	}()

	return nil
}

func (m *CommentModel) GetComments(ctx context.Context, incidentID int) ([]Comment, error) {
	query := `
		SELECT comments.*, users.name, users.role FROM comments INNER JOIN users ON comments.user_id=users.id WHERE comments.incident_id = $1 ORDER BY comments.id DESC;
	`
	rows, err := m.DB.Query(ctx, query, incidentID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %d", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Id, &comment.UserId, &comment.IncidentId, &comment.Comment, &comment.CommentUserName, &comment.CommentUserRole)
		if err != nil {
			return nil, fmt.Errorf("database query error: %v", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
