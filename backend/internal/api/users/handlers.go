package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theEricHoang/lovenote/backend/internal/api/auth"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
)

const DefaultProfilePicture = "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg"

type UserHandler struct {
	UserDAO *UserDAO
}

func NewUserHandler(userDAO *UserDAO) *UserHandler {
	return &UserHandler{UserDAO: userDAO}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	userId, err := h.UserDAO.CreateUser(r.Context(), req.Username, req.Email, profilePicture, hashedPassword)
	if err != nil {
		http.Error(w, "Error creating user in database", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(userId)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	res := struct {
		Id             uint   `json:"id"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		AccessToken    string `json:"access"`
		RefreshToken   string `json:"refresh"`
	}{
		Id:             userId,
		Username:       req.Username,
		Email:          req.Email,
		ProfilePicture: profilePicture,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding new user to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UserDAO.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusUnauthorized)
		return
	}

	err = auth.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.Id)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	res := struct {
		Id             uint   `json:"id"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		AccessToken    string `json:"access"`
		RefreshToken   string `json:"refresh"`
	}{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userId64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	userId := uint(userId64)

	user, err := h.UserDAO.GetUserById(r.Context(), userId)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	res := struct {
		Id             uint   `json:"id"`
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
		Bio            string `json:"bio"`
	}{
		Id:             user.Id,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Username       *string `json:"username,omitempty"`
		ProfilePicture *string `json:"profile_picture,omitempty"`
		Bio            *string `json:"bio,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.UserDAO.UpdateUser(r.Context(), userId, req)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully!"})
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.UserDAO.DeleteUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
