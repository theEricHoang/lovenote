package users

import (
	"encoding/json"
	"net/http"

	"github.com/theEricHoang/lovenote/backend/internal/api/auth"
)

const DefaultProfilePicture = "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg"

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		Password       string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// use default profile picture when its not provided
	profilePicture := req.ProfilePicture
	if profilePicture == "" {
		profilePicture = DefaultProfilePicture
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	userId, err := CreateUser(req.Username, req.Email, profilePicture, hashedPassword)
	if err != nil {
		http.Error(w, "Error creating user in database", http.StatusInternalServerError)
		return
	}

	res := struct {
		Id             uint   `json:"id"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
	}{
		Id:             userId,
		Username:       req.Username,
		Email:          req.Email,
		ProfilePicture: profilePicture,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding new user to JSON", http.StatusInternalServerError)
		return
	}
}
