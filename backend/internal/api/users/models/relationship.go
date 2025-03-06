package models

import (
	"encoding/json"
	"time"
)

type Relationship struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
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
