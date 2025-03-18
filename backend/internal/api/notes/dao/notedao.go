package dao

import (
	"context"

	"github.com/theEricHoang/lovenote/backend/internal/api/notes/models"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
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
