package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/handlers"
)

// define routes here
func RegisterRoutes(userHandler *handlers.UserHandler, relationshipHandler *handlers.RelationshipHandler, inviteHandler *handlers.InviteHandler) chi.Router {
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
		r.With(middleware.AuthenticateMiddleware).Patch("/me", userHandler.UpdateUserHandler)
		r.With(middleware.AuthenticateMiddleware).Delete("/me", userHandler.DeleteUserHandler)
	})

	r.Route("/api/relationships", func(r chi.Router) {
		r.With(middleware.AuthenticateMiddleware).Post("/", relationshipHandler.CreateRelationshipHandler)
		r.With(middleware.AuthenticateMiddleware).Get("/", relationshipHandler.GetUserRelationshipsHandler)
		r.Get("/{id}", relationshipHandler.GetRelationshipHandler)
		r.With(middleware.AuthenticateMiddleware).Get("/{id}/members", relationshipHandler.GetRelationshipMembersHandler)
		r.With(middleware.AuthenticateMiddleware).Patch("/{id}", relationshipHandler.UpdateRelationshipHandler)
		r.With(middleware.AuthenticateMiddleware).Delete("/{id}", relationshipHandler.DeleteRelationshipHandler)

		r.With(middleware.AuthenticateMiddleware).Post("/{id}/invite", inviteHandler.InviteUser)
	})

	r.Route("/api/invites", func(r chi.Router) {
		r.With(middleware.AuthenticateMiddleware).Get("/", inviteHandler.GetInvites)
		r.With(middleware.AuthenticateMiddleware).Post("/{id}", inviteHandler.AcceptInvite)
		r.With(middleware.AuthenticateMiddleware).Delete("/{id}", inviteHandler.DeleteInvite)
	})

	return r
}
