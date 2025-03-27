package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/theEricHoang/lovenote/backend/internal/api/auth"
	notedao "github.com/theEricHoang/lovenote/backend/internal/api/notes/dao"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
)

type contextKey string

const UserIDKey contextKey = "userID"
const RelationshipIDKey contextKey = "relationshipID"

type AuthMiddleware struct {
	AuthService *auth.AuthService
}

func NewAuthMiddleware(authService *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{AuthService: authService}
}

func (m *AuthMiddleware) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := m.AuthService.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID := claims.UserId
		expTime := claims.ExpiresAt.Time
		if time.Now().After(expTime) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type PermissionsMiddleware struct {
	RelationshipDAO *dao.RelationshipDAO
	NoteDAO         *notedao.NoteDAO
}

func NewPermissionsMiddleware(relationshipDAO *dao.RelationshipDAO, noteDAO *notedao.NoteDAO) *PermissionsMiddleware {
	return &PermissionsMiddleware{RelationshipDAO: relationshipDAO, NoteDAO: noteDAO}
}

func (m *PermissionsMiddleware) IsInRelationship(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		relationshipIDParam := chi.URLParam(r, "id")
		relationshipID64, err := strconv.ParseUint(relationshipIDParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid relationship id", http.StatusBadRequest)
			return
		}
		relationshipID := uint(relationshipID64)

		userInRelationship, err := m.RelationshipDAO.UserInRelationship(r.Context(), relationshipID, userID)
		if err != nil {
			http.Error(w, "Error checking if user is in relationship", http.StatusInternalServerError)
			return
		}
		if !userInRelationship {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RelationshipIDKey, relationshipID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *PermissionsMiddleware) IsNoteOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		noteIDParam := chi.URLParam(r, "note_id")
		noteID64, err := strconv.ParseUint(noteIDParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid note id", http.StatusBadRequest)
			return
		}
		noteID := int(noteID64)

		note, err := m.NoteDAO.GetNoteByID(r.Context(), noteID)
		if err != nil {
			http.Error(w, "Note does not exist", http.StatusBadRequest)
			return
		}

		if userID != note.Author.Id {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
