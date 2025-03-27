package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/theEricHoang/lovenote/backend/internal/api/notes/models"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"

	usermodels "github.com/theEricHoang/lovenote/backend/internal/api/users/models"
)

type NoteDAO struct {
	DB *db.Database
}

func NewNoteDAO(database *db.Database) *NoteDAO {
	return &NoteDAO{DB: database}
}

func (dao *NoteDAO) CreateNote(ctx context.Context, authorID uint, title, content, color string, x, y float32) (*models.Note, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var note models.Note
	note.Author = &usermodels.User{}
	query := `
		WITH inserted_note AS (
			INSERT INTO notes (author_id, title, content, position_x, position_y, color)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING *
		)
		SELECT
			n.id,
			a.id,
			a.username,
			a.profile_picture,
			n.title,
			n.content,
			n.position_x,
			n.position_y,
			n.color,
			n.created_at
		FROM inserted_note n
		JOIN users a ON n.author_id = a.id
	`

	row := tx.QueryRow(ctx, query, authorID, title, content, x, y, color)
	err = row.Scan(
		&note.Id,
		&note.Author.Id,
		&note.Author.Username,
		&note.Author.ProfilePicture,
		&note.Title,
		&note.Content,
		&note.PositionX,
		&note.PositionY,
		&note.Color,
		&note.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (dao *NoteDAO) GetNoteByID(ctx context.Context, noteID int) (*models.Note, error) {
	query := `
		SELECT
			n.id,
			a.id,
			a.username,
			a.profile_picture,
			n.title,
			n.content,
			n.position_x,
			n.position_y,
			n.color,
			n.created_at
		FROM notes n
		JOIN users a ON n.author_id = a.id
		WHERE n.id= $1
	`

	var note models.Note
	note.Author = &usermodels.User{}
	err := dao.DB.Pool.QueryRow(ctx, query, noteID).Scan(
		&note.Id,
		&note.Author.Id,
		&note.Author.Username,
		&note.Author.ProfilePicture,
		&note.Title,
		&note.Content,
		&note.PositionX,
		&note.PositionY,
		&note.Color,
		&note.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (dao *NoteDAO) GetNotesByRelationshipAndMonth(ctx context.Context, relationshipID, month, year int) ([]models.Note, error) {
	query := `
		SELECT
			n.id,
			a.id,
			a.username,
			a.profile_picture,
			n.title,
			n.content,
			n.position_x,
			n.position_y,
			n.color,
			n.created_at
		FROM notes n
		JOIN users a ON n.author_id = a.id
		WHERE n.relationship_id = $1
		AND EXTRACT(MONTH FROM n.created_at) = $2
		AND EXTRACT(YEAR FROM n.created_at) = $3
	`

	rows, err := dao.DB.Pool.Query(ctx, query, relationshipID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		note.Author = &usermodels.User{}
		err = rows.Scan(
			&note.Id,
			&note.Author.Id,
			&note.Author.Username,
			&note.Author.ProfilePicture,
			&note.Title,
			&note.Content,
			&note.PositionX,
			&note.PositionY,
			&note.Color,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (dao *NoteDAO) UpdateNote(ctx context.Context, noteID int, data struct {
	Title     *string  `json:"title,omitempty"`
	Content   *string  `json:"content,omitempty"`
	PositionX *float32 `json:"position_x,omitempty"`
	PositionY *float32 `json:"position_y,omitempty"`
	Color     *string  `json:"color,omitempty"`
}) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updates := []string{}
	args := []any{}
	argPos := 1

	fields := map[string]any{
		"title":      data.Title,
		"content":    data.Content,
		"position_x": data.PositionX,
		"position_y": data.PositionY,
		"color":      data.Color,
	}

	for col, val := range fields {
		if val != nil {
			updates = append(updates, fmt.Sprintf("%s = $%d", col, argPos))
			args = append(args, val)
			argPos++
		}
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE notes SET %s WHERE id = $%d", strings.Join(updates, ", "), argPos)
	args = append(args, noteID)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dao *NoteDAO) DeleteNote(ctx context.Context, noteID uint) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "DELETE FROM notes WHERE id = $1"
	_, err = tx.Exec(ctx, query, noteID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
