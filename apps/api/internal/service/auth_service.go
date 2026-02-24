package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/falasefemi2/goreact-boilerplate/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrInvalidCreds = errors.New("invalid email or password")
)

type AuthService struct {
	queries      db.Querier
	jwtSecret    string
	emailService *EmailService
}

func NewAuthService(queries db.Querier, jwtSecret string, emailService *EmailService) *AuthService {
	return &AuthService{
		queries:      queries,
		jwtSecret:    jwtSecret,
		emailService: emailService,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	// check if email is taken
	exiting, _ := s.queries.GetUserByEmail(ctx, email)
	if exiting.ID != [16]byte{} {
		return "", ErrEmailTaken
	}

	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	})
	if err != nil {
		return "", err
	}

	// Send welcome email — non-blocking
	// we don't fail registration if email fails
	go func() {
		if err := s.emailService.SendWelcome(user.Email); err != nil {
			slog.Error("failed to send welcome email",
				"email", user.Email,
				"error", err,
			)
		}
	}()

	// return jwt token
	return s.generateToken(user.ID.String())
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCreds
	}

	// Compare submitted password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCreds
	}

	return s.generateToken(user.ID.String())
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
