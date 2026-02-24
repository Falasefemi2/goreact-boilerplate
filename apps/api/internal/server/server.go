package server

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/falasefemi2/goreact-boilerplate/docs"
	"github.com/falasefemi2/goreact-boilerplate/internal/config"
	"github.com/falasefemi2/goreact-boilerplate/internal/db"
	"github.com/falasefemi2/goreact-boilerplate/internal/handler"
	appMiddleware "github.com/falasefemi2/goreact-boilerplate/internal/middleware"
	"github.com/falasefemi2/goreact-boilerplate/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"golang.org/x/time/rate"
)

func New(database *sql.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.Server.AllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	queries := db.New(database)
	emailService := service.NewEmailService(
		cfg.Email.ResendAPIKey,
		cfg.Email.FromEmail,
	)
	authService := service.NewAuthService(
		queries,
		cfg.Auth.JWTSecret,
		emailService,
	)
	authHandler := handler.NewAuthHandler(authService)
	productService := service.NewProductService(queries)
	productHandler := handler.NewProductHandler(productService)

	// Strict limiter for auth — 5 requests/minute per IP
	authLimiter := appMiddleware.NewRateLimiter(rate.Every(time.Minute/5), 5)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public auth routes with rate limiting
	r.Group(func(r chi.Router) {
		r.Use(authLimiter.Limit)
		r.Post("/api/v1/auth/register", authHandler.Register)
		r.Post("/api/v1/auth/login", authHandler.Login)
		r.Post("/api/v1/auth/logout", authHandler.Logout)
		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/docs/doc.json"),
		))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.RequireAuth(cfg.Auth.JWTSecret))
		r.Get("/api/v1/auth/me", authHandler.Me)
		r.Post("/api/v1/products", productHandler.Create)
		r.Get("/api/v1/products", productHandler.List)
		r.Get("/api/v1/products/{id}", productHandler.GetByID)
		r.Put("/api/v1/products/{id}", productHandler.Update)
		r.Delete("/api/v1/products/{id}", productHandler.Delete)
	})

	return r
}
