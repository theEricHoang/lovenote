package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	config "github.com/theEricHoang/lovenote/backend/internal"
	"github.com/theEricHoang/lovenote/backend/internal/api/auth"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
)

const DefaultProfilePicture = "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg"

var cfg config.Config = config.LoadConfig()

type UserHandler struct {
	UserDAO     *dao.UserDAO
	AuthService *auth.AuthService
}

func NewUserHandler(userDAO *dao.UserDAO, authService *auth.AuthService) *UserHandler {
	return &UserHandler{UserDAO: userDAO, AuthService: authService}
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

	hashedPassword, err := h.AuthService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user, err := h.UserDAO.CreateUser(r.Context(), req.Username, req.Email, profilePicture, hashedPassword)
	if err != nil {
		http.Error(w, "Error creating user in database", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.AuthService.GenerateTokens(r.Context(), user.Id)
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
	}{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: profilePicture,
		AccessToken:    accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusCreated)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   cfg.IsProduction,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/refresh",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

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

	err = h.AuthService.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.AuthService.GenerateTokens(r.Context(), user.Id)
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
	}{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		AccessToken:    accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   cfg.IsProduction,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/refresh",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := h.AuthService.ValidateToken(refreshToken.Value)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	realToken, err := h.AuthService.GetRefreshToken(r.Context(), claims.UserId)
	if err != nil {
		http.Error(w, "Error getting refresh token from database", http.StatusInternalServerError)
		return
	}

	if refreshToken.Value != realToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newAccess, newRefresh, err := h.AuthService.GenerateTokens(r.Context(), claims.UserId)
	if err != nil {
		http.Error(w, "Error generating new tokens", http.StatusInternalServerError)
		return
	}

	res := struct {
		AccessToken string `json:"access"`
	}{
		AccessToken: newAccess,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HttpOnly: true,
		Secure:   cfg.IsProduction,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/refresh",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	json.NewEncoder(w).Encode(res)
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
		return
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

func (h *UserHandler) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing 'username' query parameter", http.StatusBadRequest)
		return
	}

	// Default values for pagination
	limit := 10
	page := 1

	// Parse limit and page (handle errors gracefully)
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	offset := (page - 1) * limit

	users, userCount, err := h.UserDAO.SearchUsersByName(r.Context(), username, limit, offset)
	if err != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	baseUrl := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
	queryParams := fmt.Sprintf("username=%s&limit=%d", username, limit)

	var nextLink, prevLink *string
	if offset+limit < userCount {
		next := fmt.Sprintf("%s?%s&page=%d", baseUrl, queryParams, page+1)
		nextLink = &next
	}
	if page > 1 {
		prev := fmt.Sprintf("%s?%s&page=%d", baseUrl, queryParams, page-1)
		prevLink = &prev
	}

	response := map[string]any{
		"count": userCount,
		"next":  nextLink,
		"prev":  prevLink,
		"users": users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
