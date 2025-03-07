package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/theEricHoang/lovenote/backend/internal/api/users/models"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
)

type RelationshipDAO struct {
	DB *db.Database
}

func NewRelationshipDAO(database *db.Database) *RelationshipDAO {
	return &RelationshipDAO{DB: database}
}

func (dao *RelationshipDAO) CreateRelationship(ctx context.Context, name, picture string) (*models.Relationship, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var relationship models.Relationship
	query := "INSERT INTO relationships (name, picture) values ($1, $2) RETURNING id, name, picture, created_at"

	row := tx.QueryRow(ctx, query, name, picture)
	err = row.Scan(&relationship.Id, &relationship.Name, &relationship.Picture, &relationship.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &relationship, nil
}

func (dao *RelationshipDAO) GetRelationshipById(ctx context.Context, id uint) (*models.Relationship, error) {
	var relationship models.Relationship
	query := "SELECT id, name, picture, created_at FROM relationships WHERE id = $1"

	row := dao.DB.Pool.QueryRow(ctx, query, id)
	err := row.Scan(&relationship.Id, &relationship.Name, &relationship.Picture, &relationship.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
}

func (dao *RelationshipDAO) UpdateRelationship(ctx context.Context, relationshipId uint, data struct {
	Name    *string `json:"name,omitempty"`
	Picture *string `json:"picture,omitempty"`
}) (*models.Relationship, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var relationship models.Relationship
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	fields := map[string]*string{
		"name":    data.Name,
		"picture": data.Picture,
	}

	for col, val := range fields {
		if val != nil {
			updates = append(updates, fmt.Sprintf("%s = %d", col, argPos))
			args = append(args, *val)
			argPos++
		}
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	query := fmt.Sprintf("UPDATE relationships SET %s WHERE id = $%d RETURNING id, name, picture, created_at", strings.Join(updates, ", "), argPos)
	args = append(args, relationshipId)

	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&relationship.Id, &relationship.Name, &relationship.Picture, &relationship.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &relationship, nil
}

func (dao *RelationshipDAO) DeleteRelationship(ctx context.Context, id uint) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "DELETE FROM relationships WHERE id = $1"
	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RelationshipDAO) UserInRelationship(ctx context.Context, relationshipId, userId uint) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM relationship_members WHERE relationship_id = $1 AND user_id = $2
	)`

	err := dao.DB.Pool.QueryRow(ctx, query, relationshipId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dao *RelationshipDAO) IsUserOnlyMember(ctx context.Context, userID, relationshipID uint) (bool, error) {
	var isOnly bool
	query := `SELECT COUNT(*) = 1 FROM relationship_members WHERE relationship_id = $1 AND user_id = $2`

	err := dao.DB.Pool.QueryRow(ctx, query, relationshipID, userID).Scan(&isOnly)
	if err != nil {
		return false, err
	}
	return isOnly, nil
}

func (dao *RelationshipDAO) AddUserToRelationship(ctx context.Context, userID, relationshipID uint) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO relationship_members (relationship_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	_, err = tx.Exec(ctx, query, relationshipID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
