package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo      Repository
	jwtSecret string
	log       logger.Logger
}

func NewService(repo Repository, jwtSecret string, log logger.Logger) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (s *service) Register(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	user := &User{
		ID:           newID(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := generateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) GetByID(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	user, err := s.repo.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func newID() string {
	// 16 bytes => 32 hex chars
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to time-based ID if rand fails
		return hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(b)
}
