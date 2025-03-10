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

func (dao *RelationshipDAO) CreateRelationshipAndAddUser(ctx context.Context, name, picture string, userID uint) (*models.Relationship, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// count to ensure user is not in more than 10 relationships
	var count int
	countQuery := "SELECT COUNT(*) FROM relationship_members WHERE user_id = $1"
	err = tx.QueryRow(ctx, countQuery, userID).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count >= 10 {
		return nil, fmt.Errorf("user is in maximum relationships (10)")
	}

	// create relationship
	var relationship models.Relationship
	query := "INSERT INTO relationships (name, picture) values ($1, $2) RETURNING id, name, picture, created_at"

	row := tx.QueryRow(ctx, query, name, picture)
	err = row.Scan(&relationship.Id, &relationship.Name, &relationship.Picture, &relationship.CreatedAt)
	if err != nil {
		return nil, err
	}

	// add user to new relationship
	query = `INSERT INTO relationship_members (relationship_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	_, err = tx.Exec(ctx, query, relationship.Id, userID)
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

	// count to ensure user is not in more than 10 relationships
	var count int
	countQuery := "SELECT COUNT(*) FROM relationship_members WHERE user_id = $1"
	err = tx.QueryRow(ctx, countQuery, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= 10 {
		return fmt.Errorf("user is in maximum relationships (10)")
	}

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

func (dao *RelationshipDAO) GetUserRelationships(ctx context.Context, userID uint) ([]models.Relationship, error) {
	query := `
		SELECT r.id, r.name, r.picture, r.created_at
		FROM relationships r
		INNER JOIN relationship_members rm ON r.id = rm.relationship_id
		WHERE rm.user_id = $1
	`

	rows, err := dao.DB.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []models.Relationship
	for rows.Next() {
		var relationship models.Relationship
		if err := rows.Scan(&relationship.Id, &relationship.Name, &relationship.Picture, &relationship.CreatedAt); err != nil {
			return nil, err
		}
		relationships = append(relationships, relationship)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relationships, nil
}
