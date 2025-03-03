CREATE INDEX idx_invites_relationship_id ON invites(relationship_id);
CREATE INDEX idx_invites_inviter_id ON invites(inviter_id);
CREATE INDEX idx_invites_invitee_id ON invites(invitee_id);

CREATE INDEX idx_relationship_members_relationship_id ON relationship_members(relationship_id);
CREATE INDEX idx_relationship_members_user_id ON relationship_members(user_id);

CREATE INDEX idx_relationship_notes_relationship_id ON relationship_notes(relationship_id);
CREATE INDEX idx_relationship_notes_note_id ON relationship_notes(note_id);

CREATE INDEX idx_notes_author_id ON notes(author_id);