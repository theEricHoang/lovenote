package dao

import (
	"context"

	"github.com/theEricHoang/lovenote/backend/internal/api/users/models"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
)

type InviteDAO struct {
	DB *db.Database
}

func NewInviteDAO(database *db.Database) *InviteDAO {
	return &InviteDAO{DB: database}
}

func (dao *InviteDAO) CreateInvite(ctx context.Context, relationshipId, inviterId, inviteeId uint, body string) (*models.Invite, error) {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var invite models.Invite
	query := `WITH inserted_invite AS (
		INSERT INTO invites (relationship_id, inviter_id, invitee_id, body)
		VALUES ($1, $2, $3, $4)
		RETURNING *
		)
		SELECT
			ii.id,
			r.id,
			r.name,
			r.picture,
			inviter.id,
			inviter.username,
			inviter.profile_picture,
			invitee.id,
			invitee.username,
			invitee.profile_picture,
			ii.body
		FROM inserted_invite ii
		JOIN relationships r ON ii.relationship_id = r.id
		JOIN users inviter ON ii.inviter_id = inviter.id
		JOIN users invitee ON ii.invitee_id = invitee.id;`

	row := tx.QueryRow(ctx, query, relationshipId, inviterId, inviteeId, body)
	err = row.Scan(
		&invite.Id,
		&invite.Relationship.Id,
		&invite.Relationship.Name,
		&invite.Relationship.Picture,
		&invite.Inviter.Id,
		&invite.Inviter.Username,
		&invite.Inviter.ProfilePicture,
		&invite.Invitee.Id,
		&invite.Invitee.Username,
		&invite.Invitee.ProfilePicture,
		&invite.Body,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

func (dao *InviteDAO) GetInviteById(ctx context.Context, inviteId uint) (*models.Invite, error) {
	var invite models.Invite
	query := "SELECT id, relationship_id, inviter_id, invitee_id, body FROM invites WHERE id = $1"
	row := dao.DB.Pool.QueryRow(ctx, query, inviteId)
	err := row.Scan(&invite.Id, &invite.Relationship.Id, &invite.Inviter.Id, &invite.Invitee.Id, &invite.Body)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (dao *InviteDAO) DeleteInvite(ctx context.Context, inviteId uint) error {
	tx, err := dao.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "DELETE FROM invites WHERE id = $1"
	_, err = tx.Exec(ctx, query, inviteId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
