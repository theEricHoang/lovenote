ALTER TABLE invites ADD CONSTRAINT unique_invite UNIQUE (relationship_id, invitee_id);
