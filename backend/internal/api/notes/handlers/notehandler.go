package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
	note, err := h.NoteDAO.CreateNote(r.Context(), authorID, relationshipID, req.Title, req.Content, req.Color, req.PositionX, req.PositionY)
	if err != nil {
		http.Error(w, "Error inserting note into database", http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) GetRelationshipNotes(w http.ResponseWriter, r *http.Request) {
	relationshipID, ok := r.Context().Value(middleware.RelationshipIDKey).(uint)
	if !ok {
		http.Error(w, "Missing relationship ID", http.StatusBadRequest)
		return
	}

	// default month and year
	now := time.Now()
	month := int(now.Month())
	year := int(now.Year())

	// parse month and year
	if m, err := strconv.Atoi(r.URL.Query().Get("month")); err == nil && m > 0 && m < 13 {
		month = m
	}
	if y, err := strconv.Atoi(r.URL.Query().Get("year")); err == nil && y > 0 {
		year = y
	}

	notes, err := h.NoteDAO.GetNotesByRelationshipAndMonth(r.Context(), relationshipID, month, year)
	if err != nil {
		http.Error(w, "Error getting notes from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) EditNote(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "note_id")
	noteID64, err := strconv.ParseUint(noteIDParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	noteID := int(noteID64)

	var req struct {
		Title     *string  `json:"title,omitempty"`
		Content   *string  `json:"content,omitempty"`
		PositionX *float32 `json:"position_x,omitempty"`
		PositionY *float32 `json:"position_y,omitempty"`
		Color     *string  `json:"color,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.NoteDAO.UpdateNote(r.Context(), noteID, req)
	if err != nil {
		http.Error(w, "Error updating note in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Note updated successfully"})
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "note_id")
	noteID64, err := strconv.ParseUint(noteIDParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid relationship id", http.StatusBadRequest)
		return
	}
	noteID := uint(noteID64)

	err = h.NoteDAO.DeleteNote(r.Context(), noteID)
	if err != nil {
		http.Error(w, "Error deleting note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
