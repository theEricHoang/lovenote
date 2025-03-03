package models

import "encoding/json"

type Invite struct {
	Id           uint         `json:"id"`
	Relationship Relationship `json:"relationship"`
	Inviter      User         `json:"inviter"`
	InviteeId    uint         `json:"invitee_id"`
	Body         string       `json:"body"`
}

func (i *Invite) ToJSON(view string) ([]byte, error) {
	switch view {
	case "minimal":
		return json.Marshal(struct {
			Id           uint         `json:"id"`
			Relationship Relationship `json:"relationship"`
			Inviter      User         `json:"inviter"`
		}{
			Id:           i.Id,
			Relationship: i.Relationship,
			Inviter:      i.Inviter,
		})
	case "full":
		return json.Marshal(i)
	default:
		return json.Marshal(i)
	}
}
