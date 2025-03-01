package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
)

const DefaultRelationshipPicture = "https://img.freepik.com/free-vector/gradient-heart_78370-478.jpg"

type RelationshipHandler struct {
	RelationshipDAO *dao.RelationshipDAO
}

func NewRelationshipHandler(relationshipDAO *dao.RelationshipDAO) *RelationshipHandler {
	return &RelationshipHandler{RelationshipDAO: relationshipDAO}
}

func (h *RelationshipHandler) CreateRelationshipHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	picture := req.Picture
	if picture == "" {
		picture = DefaultRelationshipPicture
	}

	relationship, err := h.RelationshipDAO.CreateRelationship(r.Context(), req.Name, picture)
	if err != nil {
		http.Error(w, "Error creating relationship in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(relationship)
	if err != nil {
		http.Error(w, "Error encoding new user to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *RelationshipHandler) GetRelationshipHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	relationshipId64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	relationshipId := uint(relationshipId64)

	relationship, err := h.RelationshipDAO.GetRelationshipById(r.Context(), relationshipId)
	if err != nil {
		http.Error(w, "Relationship does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(relationship)
	if err != nil {
		http.Error(w, "Error encoding relationship to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *RelationshipHandler) UpdateRelationshipHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	id := uint(id64)

	userInRelationship, err := h.RelationshipDAO.UserInRelationship(r.Context(), id, userId)
	if err != nil {
		http.Error(w, "Error checking if current user is in the requested relationship", http.StatusInternalServerError)
		return
	}
	if !userInRelationship {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name    *string `json:"name,omitempty"`
		Picture *string `json:"picture,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	relationship, err := h.RelationshipDAO.UpdateRelationship(r.Context(), id, req)
	if err != nil {
		http.Error(w, "Error updating relationship", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(relationship)
	if err != nil {
		http.Error(w, "Error encoding relationship to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *RelationshipHandler) DeleteRelationshipHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	id := uint(id64)

	// check to see if user is the only person in relationship
	isOnly, err := h.RelationshipDAO.IsUserOnlyMember(r.Context(), userId, id)
	if err != nil {
		http.Error(w, "Error checking permissions", http.StatusInternalServerError)
		return
	}
	if !isOnly {
		http.Error(w, "Unauthorized, relationships can only be deleted if only one person belongs to them", http.StatusUnauthorized)
		return
	}

	err = h.RelationshipDAO.DeleteRelationship(r.Context(), id)
	if err != nil {
		http.Error(w, "Error deleting relationship", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
