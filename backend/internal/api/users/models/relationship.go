package models

import (
	"encoding/json"
	"time"
)

type Relationship struct {
	Id        uint       `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Picture   string     `json:"picture,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (r *Relationship) ToJSON(view string) ([]byte, error) {
	switch view {
	case "minimal":
		return json.Marshal(struct {
			Id      uint   `json:"id"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}{
			Id:      r.Id,
			Name:    r.Name,
			Picture: r.Picture,
		})
	case "full":
		return json.Marshal(r)
	default:
		return json.Marshal(r)
	}
}
