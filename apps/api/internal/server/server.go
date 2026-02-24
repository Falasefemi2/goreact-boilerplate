package server

import (
	"database/sql"
	"net/http"

	"github.com/falasefemi2/goreact-boilerplate/internal/config"
	"github.com/falasefemi2/goreact-boilerplate/internal/db"
	"github.com/falasefemi2/goreact-boilerplate/internal/handler"
	appMiddleware "github.com/falasefemi2/goreact-boilerplate/internal/middleware"
	"github.com/falasefemi2/goreact-boilerplate/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(database *sql.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigin}, // must be exact, not *
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // required for cookies
		MaxAge:           300,
	}))

	queries := db.New(database)

	authService := service.NewAuthService(queries, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	productService := service.NewProductService(queries)
	productHandler := handler.NewProductHandler(productService)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.RequireAuth(cfg.JWTSecret))

		// Auth
		r.Get("/api/v1/auth/me", authHandler.Me)

		// Products
		r.Post("/api/v1/products", productHandler.Create)
		r.Get("/api/v1/products", productHandler.List)
		r.Get("/api/v1/products/{id}", productHandler.GetByID)
		r.Put("/api/v1/products/{id}", productHandler.Update)
		r.Delete("/api/v1/products/{id}", productHandler.Delete)
	})

	// Protected routes — anything here requires a valid JWT
	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/logout", authHandler.Logout)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
