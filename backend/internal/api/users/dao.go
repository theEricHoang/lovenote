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

func GetUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, email, profile_picture, password_hash FROM users WHERE username = $1"
	err := db.DB.QueryRow(context.Background(), query, username).Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
