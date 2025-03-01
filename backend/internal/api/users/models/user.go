package models

import "encoding/json"

type User struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
	PasswordHash   string `json:"-"`
	CreatedAt      string `json:"created_at"`
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
