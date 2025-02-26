package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users"
)

// define routes here
func RegisterRoutes(userHandler *users.UserHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.StripSlashes)
	r.Use(chimiddleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// users routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.RegisterHandler)
		r.Post("/login", userHandler.LoginHandler)
		r.Get("/{id}", userHandler.GetUserHandler)
	})

	return r
}
