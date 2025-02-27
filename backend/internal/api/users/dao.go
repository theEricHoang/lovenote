package users

import (
	"context"
	"fmt"
	"strings"

	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
)

type UserDAO struct {
	DB *db.Database
}

func NewUserDAO(database *db.Database) *UserDAO {
	return &UserDAO{DB: database}
}

func (dao *UserDAO) CreateUser(ctx context.Context, username, email, profilePicture, passwordHash string) (uint, error) {
	var userId uint
	query := `INSERT INTO users (username, email, profile_picture, password_hash) values ($1, $2, $3, $4) RETURNING id`
	err := dao.DB.Pool.QueryRow(ctx, query, username, email, profilePicture, passwordHash).Scan(&userId)
	return userId, err
}

func (dao *UserDAO) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := "SELECT id, username, email, profile_picture, bio, password_hash FROM users WHERE username = $1"
	err := dao.DB.Pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.Bio, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserById(ctx context.Context, id uint) (*User, error) {
	var user User
	query := "SELECT id, username, email, profile_picture, bio, password_hash FROM users WHERE id = $1"
	row := dao.DB.Pool.QueryRow(ctx, query, id)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.Bio, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) UpdateUser(ctx context.Context, userId uint, data struct {
	Username       *string `json:"username,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
}) error {
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	fields := map[string]*string{
		"username":        data.Username,
		"profile_picture": data.ProfilePicture,
		"bio":             data.Bio,
	}

	for col, val := range fields {
		if val != nil {
			updates = append(updates, fmt.Sprintf("%s = $%d", col, argPos))
			args = append(args, *val)
			argPos++
		}
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE USERS SET %s WHERE id = $%d", strings.Join(updates, ", "), argPos)
	args = append(args, userId)

	_, err := dao.DB.Pool.Exec(ctx, query, args...)
	return err
}

func (dao *UserDAO) DeleteUser(ctx context.Context, userId uint) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := dao.DB.Pool.Exec(ctx, query, userId)
	return err
}
