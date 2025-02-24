package users

import (
	"context"

	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
)

func CreateUser(username, email, profilePicture, passwordHash string) (uint, error) {
	var userId uint
	query := `INSERT INTO users (username, email, profile_picture, password_hash) values ($1, $2, $3, $4) RETURNING id`
	err := db.DB.QueryRow(context.Background(), query, username, email, profilePicture, passwordHash).Scan(&userId)
	return userId, err
}
