package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id             uint       `json:"id,omitempty"`
	Username       string     `json:"username,omitempty"`
	Email          string     `json:"email,omitempty"`
	ProfilePicture string     `json:"profile_picture,omitempty"`
	Bio            string     `json:"bio,omitempty"`
	PasswordHash   string     `json:"-"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

func (u *User) ToJSON(view string) ([]byte, error) {
	switch view {
	case "minimal":
		return json.Marshal(struct {
			Id             uint   `json:"id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
		}{
			Id:             u.Id,
			Username:       u.Username,
			ProfilePicture: u.ProfilePicture,
		})
	case "profile":
		return json.Marshal(struct {
			Id             uint   `json:"id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
			Bio            string `json:"bio"`
		}{
			Id:             u.Id,
			Username:       u.Username,
			ProfilePicture: u.ProfilePicture,
			Bio:            u.Bio,
		})
	case "full":
		return json.Marshal(u)
	default:
		return json.Marshal(u)
	}
}
