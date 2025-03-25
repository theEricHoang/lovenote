package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/notes/dao"
	usersdao "github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
)

type NoteHandler struct {
	NoteDAO         *dao.NoteDAO
	RelationshipDAO *usersdao.RelationshipDAO
}

func NewNoteHandler(noteDAO *dao.NoteDAO, relationshipDAO *usersdao.RelationshipDAO) *NoteHandler {
	return &NoteHandler{NoteDAO: noteDAO, RelationshipDAO: relationshipDAO}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	// get author info
	authorID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get relationship ID from url
	relationshipIDParam := chi.URLParam(r, "id")
	relationshipID64, err := strconv.ParseUint(relationshipIDParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	relationshipID := uint(relationshipID64)

	// check if author is in relationship
	authorInRelationship, err := h.RelationshipDAO.UserInRelationship(r.Context(), relationshipID, authorID)
	if err != nil {
		http.Error(w, "Error verifying if user is in relationship", http.StatusInternalServerError)
		return
	}
	if !authorInRelationship {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title     string  `json:"title"`
		Content   string  `json:"content"`
		PositionX float32 `json:"position_x"`
		PositionY float32 `json:"position_y"`
		Color     string  `json:"color"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// input validation checks
	if len(req.Content) > 500 {
		http.Error(w, "Content too long. Max is 500", http.StatusBadRequest)
		return
	}
	if len(req.Title) > 100 {
		http.Error(w, "Title too long. Max is 100", http.StatusBadRequest)
		return
	}

	// create new note
	note, err := h.NoteDAO.CreateNote(r.Context(), authorID, req.Title, req.Content, req.Color, req.PositionX, req.PositionY)
	if err != nil {
		http.Error(w, "Error inserting note into database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}
