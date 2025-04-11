package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/theEricHoang/lovenote/backend/internal/api/users/models"
	"github.com/theEricHoang/lovenote/backend/internal/pkg/db"
)

type UserDAO struct {
	DB *db.Database
}

func NewUserDAO(database *db.Database) *UserDAO {
	return &UserDAO{DB: database}
}

func (dao *UserDAO) CreateUser(ctx context.Context, username, email, profilePicture, passwordHash string) (*models.User, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var user models.User
	query := `INSERT INTO users (username, email, profile_picture, password_hash) values ($1, $2, $3, $4)
		RETURNING id, username, email, profile_picture, password_hash`
	err = tx.QueryRow(ctx, query, username, email, profilePicture, passwordHash).Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dao *UserDAO) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, profile_picture, bio, password_hash FROM users WHERE username = $1"
	err := dao.DB.Pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.Bio, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
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
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updates := []string{}
	args := []any{}
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

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updates, ", "), argPos)
	args = append(args, userId)

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

func (dao *UserDAO) DeleteUser(ctx context.Context, userId uint) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "DELETE FROM users WHERE id = $1"
	_, err = tx.Exec(ctx, query, userId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDAO) SearchUsersByName(ctx context.Context, search string, limit, offset int) ([]models.User, int, error) {
	query := `
		SELECT id, username, email, profile_picture, bio
		FROM users
		WHERE username ILIKE '%' || $1 || '%'
		ORDER BY username
		LIMIT $2 OFFSET $3
	`
	// limit: how many results to return per page
	// offset: how many results to skip

	rows, err := dao.DB.Pool.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.ProfilePicture, &user.Bio); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM users WHERE USERNAME ILIKE '%' || $1 || '%'`
	err = dao.DB.Pool.QueryRow(ctx, countQuery, search).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
