package notes

import (
	"encoding/json"

	"github.com/theEricHoang/lovenote/backend/internal/api/users"
)

type Note struct {
	Id        uint       `json:"id"`
	Author    users.User `json:"author"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	PositionX float32    `json:"position_x"`
	PositionY float32    `json:"position_y"`
	CreatedAt string     `json:"created_at"`
}

func (n *Note) ToJSON(view string) ([]byte, error) {
	return json.Marshal(n)
}
