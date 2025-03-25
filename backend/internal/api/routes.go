package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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
) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.StripSlashes)
	r.Use(chimiddleware.Logger)

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
		r.With(authMiddleware.AuthenticateMiddleware).Patch("/me", userHandler.UpdateUserHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Delete("/me", userHandler.DeleteUserHandler)
	})

	r.Route("/api/relationships", func(r chi.Router) {
		r.With(authMiddleware.AuthenticateMiddleware).Post("/", relationshipHandler.CreateRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Get("/", relationshipHandler.GetUserRelationshipsHandler)
		r.Get("/{id}", relationshipHandler.GetRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Get("/{id}/members", relationshipHandler.GetRelationshipMembersHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Patch("/{id}", relationshipHandler.UpdateRelationshipHandler)
		r.With(authMiddleware.AuthenticateMiddleware).Delete("/{id}", relationshipHandler.DeleteRelationshipHandler)

		r.With(authMiddleware.AuthenticateMiddleware).Post("/{id}", noteHandler.CreateNote)

		r.With(authMiddleware.AuthenticateMiddleware).Post("/{id}/invite", inviteHandler.InviteUser)
	})

	r.Route("/api/invites", func(r chi.Router) {
		r.With(authMiddleware.AuthenticateMiddleware).Get("/", inviteHandler.GetInvites)
		r.With(authMiddleware.AuthenticateMiddleware).Post("/{id}", inviteHandler.AcceptInvite)
		r.With(authMiddleware.AuthenticateMiddleware).Delete("/{id}", inviteHandler.DeleteInvite)
	})

	return r
}
