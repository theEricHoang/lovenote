package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	notehandlers "github.com/theEricHoang/lovenote/backend/internal/api/notes/handlers"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/handlers"
)

// define routes here
func RegisterRoutes(
	userHandler *handlers.UserHandler,
	relationshipHandler *handlers.RelationshipHandler,
	inviteHandler *handlers.InviteHandler,
	noteHandler *notehandlers.NoteHandler,
	authMiddleware *middleware.AuthMiddleware,
	permissionsMiddleware *middleware.PermissionsMiddleware,
) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.StripSlashes)
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Allow cookies/auth headers
		MaxAge:           300,  // Cache CORS response for 5 minutes
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// users routes
	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", userHandler.SearchUsersHandler)
		r.Post("/", userHandler.RegisterHandler)
		r.Post("/login", userHandler.LoginHandler)
		r.Get("/{id}", userHandler.GetUserHandler)
		r.Post("/refresh", userHandler.RefreshTokenHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Get("/me", userHandler.GetSelfHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Patch("/me", userHandler.UpdateUserHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Delete("/me", userHandler.DeleteUserHandler)
	})

	r.Route("/api/relationships", func(r chi.Router) {
		r.With(authMiddleware.AuthenticateMiddleware).Post("/", relationshipHandler.CreateRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Get("/", relationshipHandler.GetUserRelationshipsHandler)
		r.Get("/{id}", relationshipHandler.GetRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Get("/{id}/members", relationshipHandler.GetRelationshipMembersHandler)
		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship).Patch("/{id}", relationshipHandler.UpdateRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship).Delete("/{id}", relationshipHandler.DeleteRelationshipHandler)

		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship).Post("/{id}/notes", noteHandler.CreateNote)
		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship).Get("/{id}/notes", noteHandler.GetRelationshipNotes)
		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship, permissionsMiddleware.IsNoteOwner).Patch("/{id}/notes/{note_id}", noteHandler.EditNote)
		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship, permissionsMiddleware.IsNoteOwner).Delete("/{id}/notes/{note_id}", noteHandler.DeleteNote)

		r.With(authMiddleware.AuthenticateMiddleware, permissionsMiddleware.IsInRelationship).Post("/{id}/invite", inviteHandler.InviteUser)
	})

	r.Route("/api/invites", func(r chi.Router) {
		r.With(authMiddleware.AuthenticateMiddleware).Get("/", inviteHandler.GetInvites)
		r.With(authMiddleware.AuthenticateMiddleware).Post("/{id}", inviteHandler.AcceptInvite)
		r.With(authMiddleware.AuthenticateMiddleware).Delete("/{id}", inviteHandler.DeleteInvite)
	})

	return r
}
