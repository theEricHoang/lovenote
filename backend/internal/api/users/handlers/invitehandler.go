package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
)

type InviteHandler struct {
	InviteDAO       *dao.InviteDAO
	RelationshipDAO *dao.RelationshipDAO
}

func NewInviteHandler(inviteDAO *dao.InviteDAO, relationshipDAO *dao.RelationshipDAO) *InviteHandler {
	return &InviteHandler{InviteDAO: inviteDAO, RelationshipDAO: relationshipDAO}
}

func (h *InviteHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	// get inviter info
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get relationship ID from url
	relationshipIdParam := chi.URLParam(r, "id")
	relationshipId64, err := strconv.ParseUint(relationshipIdParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	relationshipId := uint(relationshipId64)

	// check if inviter is in relationship
	inviterInRelationship, err := h.RelationshipDAO.UserInRelationship(r.Context(), relationshipId, userId)
	if err != nil {
		http.Error(w, "Error checking if user is in relationship", http.StatusInternalServerError)
		return
	}
	if !inviterInRelationship {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		InviteeId uint   `json:"invitee_id"`
		Body      string `json:"body"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Body) > 255 {
		http.Error(w, "Body too long. Max is 255", http.StatusBadRequest)
		return
	}

	// TODO: Check if invite already exists before creating new invite

	// create new invite
	invite, err := h.InviteDAO.CreateInvite(r.Context(), relationshipId, userId, req.InviteeId, req.Body)
	if err != nil {
		http.Error(w, "Error inserting invite into database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(invite)
	if err != nil {
		http.Error(w, "Error encoding new invite to JSON", http.StatusInternalServerError)
		return
	}
}

func (h *InviteHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	// get current user
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get invite id from url
	inviteIdParam := chi.URLParam(r, "id")
	inviteId64, err := strconv.ParseUint(inviteIdParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	inviteId := uint(inviteId64)

	// confirm if current user is invitee
	invite, err := h.InviteDAO.GetInviteById(r.Context(), inviteId)
	if err != nil {
		http.Error(w, "Error checking if user is invitee", http.StatusInternalServerError)
		return
	}
	if userId != invite.Invitee.Id {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// add them to the relationship
	err = h.RelationshipDAO.AddUserToRelationship(r.Context(), userId, invite.Relationship.Id)
	if err != nil {
		http.Error(w, "Error adding user to relationship", http.StatusInternalServerError)
		return
	}

	// delete invite
	err = h.InviteDAO.DeleteInvite(r.Context(), inviteId)
	if err != nil {
		http.Error(w, "Error deleting invite", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message":         "User added to relationship",
		"relationship_id": strconv.FormatUint(uint64(invite.Relationship.Id), 10),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "error writing response to json", http.StatusInternalServerError)
		return
	}
}

func (h *InviteHandler) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	// get current user
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get invite id from url
	inviteIdParam := chi.URLParam(r, "id")
	inviteId64, err := strconv.ParseUint(inviteIdParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	inviteId := uint(inviteId64)

	// confirm if current user is invitee or invitee
	invite, err := h.InviteDAO.GetInviteById(r.Context(), inviteId)
	if err != nil {
		http.Error(w, "Error checking if user is invitee", http.StatusInternalServerError)
		return
	}
	if userId != invite.Invitee.Id || userId != invite.Inviter.Id {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// delete invite
	err = h.InviteDAO.DeleteInvite(r.Context(), inviteId)
	if err != nil {
		http.Error(w, "Error deleting invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
