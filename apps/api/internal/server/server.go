package server

import (
	"database/sql"
	"net/http"

	"github.com/falasefemi2/goreact-boilerplate/internal/db"
	"github.com/falasefemi2/goreact-boilerplate/internal/handler"
	appMiddleware "github.com/falasefemi2/goreact-boilerplate/internal/middleware"
	"github.com/falasefemi2/goreact-boilerplate/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(database *sql.DB, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	queries := db.New(database)

	authService := service.NewAuthService(queries, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	})

	// Protected routes — anything here requires a valid JWT
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.RequireAuth(jwtSecret))
		r.Get("/auth/me", authHandler.Me)
	})

	return r
}
