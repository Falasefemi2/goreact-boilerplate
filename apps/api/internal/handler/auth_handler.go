package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/falasefemi2/goreact-boilerplate/internal/middleware"
	"github.com/falasefemi2/goreact-boilerplate/internal/response"
	"github.com/falasefemi2/goreact-boilerplate/internal/service"
	appvalidator "github.com/falasefemi2/goreact-boilerplate/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type authRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type authResponse struct {
	Token string `json:"token"`
}

// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authRequest true "Register request"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	token, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			response.Error(w, http.StatusConflict, "email already in use")
			return
		}
		response.Error(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	setAuthCookie(w, token)
	response.JSON(w, http.StatusCreated, map[string]string{"message": "registered successfully"})
}

// @Summary      Login
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authRequest true "Login request"
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			response.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		response.Error(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	setAuthCookie(w, token)
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged in successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // delete the cookie immediately
	})
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// setAuthCookie sets the JWT as an httpOnly cookie
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true, // not accessible via JavaScript
		Path:     "/",
		MaxAge:   86400, // 24 hours in seconds
		SameSite: http.SameSiteLaxMode,
		// Secure: true  // uncomment in production (requires HTTPS)
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	response.JSON(w, http.StatusOK, map[string]string{"user_id": userID})
}
