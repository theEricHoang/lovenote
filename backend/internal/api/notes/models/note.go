package models

import (
	"encoding/json"

	"github.com/theEricHoang/lovenote/backend/internal/api/users/models"
)

type Note struct {
	Id        uint        `json:"id"`
	Author    models.User `json:"author"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	PositionX float32     `json:"position_x"`
	PositionY float32     `json:"position_y"`
	Color     string      `json:"color"`
	CreatedAt string      `json:"created_at"`
}

func (n *Note) ToJSON(view string) ([]byte, error) {
	return json.Marshal(n)
}
