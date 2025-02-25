package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users"
)

// define routes here
func RegisterRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.StripSlashes)
	r.Use(chimiddleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// users routes
	r.Post("/users", users.RegisterHandler)
	r.Post("/users/login", users.LoginHandler)

	return r
}
