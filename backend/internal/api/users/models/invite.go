package models

import "encoding/json"

type Invite struct {
	Id           uint          `json:"id,omitempty"`
	Relationship *Relationship `json:"relationship,omitempty"`
	Inviter      *User         `json:"inviter,omitempty"`
	Invitee      *User         `json:"invitee,omitempty"`
	Body         string        `json:"body,omitempty"`
}

func (i *Invite) ToJSON(view string) ([]byte, error) {
	switch view {
	case "minimal":
		return json.Marshal(struct {
			Id           uint          `json:"id"`
			Relationship *Relationship `json:"relationship"`
			Inviter      *User         `json:"inviter"`
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
