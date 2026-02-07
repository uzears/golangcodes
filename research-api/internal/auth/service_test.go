package auth

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	created       *User
	findByEmail   func(email string) (*User, error)
	findByID      func(id string) (*User, error)
	createErr     error
	findEmailErr  error
	findByIDErr   error
}

func (r *fakeRepo) Create(ctx context.Context, user *User) error {
	r.created = user
	return r.createErr
}

func (r *fakeRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	if r.findByEmail != nil {
		return r.findByEmail(email)
	}
	return nil, r.findEmailErr
}

func (r *fakeRepo) FindById(ctx context.Context, userID string) (*User, error) {
	if r.findByID != nil {
		return r.findByID(userID)
	}
	return nil, r.findByIDErr
}

func TestRegisterSuccess(t *testing.T) {
	repo := &fakeRepo{
		findByEmail: func(email string) (*User, error) {
			return nil, nil
		},
	}
	svc := NewService(repo, "secret", nil)

	userID, err := svc.Register(context.Background(), "a@b.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID == "" {
		t.Fatalf("expected userID to be set")
	}
	if repo.created == nil {
		t.Fatalf("expected user to be created")
	}
	if repo.created.Email != "a@b.com" {
		t.Fatalf("unexpected email: %s", repo.created.Email)
	}
	if repo.created.PasswordHash == "password123" {
		t.Fatalf("password should be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.created.PasswordHash), []byte("password123")); err != nil {
		t.Fatalf("stored password hash is invalid")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	repo := &fakeRepo{
		findByEmail: func(email string) (*User, error) {
			return &User{ID: "u1", Email: email}, nil
		},
	}
	svc := NewService(repo, "secret", nil)

	if _, err := svc.Register(context.Background(), "a@b.com", "password123"); err == nil {
		t.Fatalf("expected error for duplicate email")
	}
}

func TestLoginSuccess(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &fakeRepo{
		findByEmail: func(email string) (*User, error) {
			return &User{
				ID:           "u1",
				Email:        email,
				PasswordHash: string(hash),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}
	svc := NewService(repo, "secret", nil)

	token, err := svc.Login(context.Background(), "a@b.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token to be returned")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &fakeRepo{
		findByEmail: func(email string) (*User, error) {
			return &User{
				ID:           "u1",
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}
	svc := NewService(repo, "secret", nil)

	if _, err := svc.Login(context.Background(), "a@b.com", "wrong"); err == nil {
		t.Fatalf("expected error for invalid credentials")
	}
}

func TestGetByIDSuccess(t *testing.T) {
	repo := &fakeRepo{
		findByID: func(id string) (*User, error) {
			return &User{ID: id, Email: "a@b.com"}, nil
		},
	}
	svc := NewService(repo, "secret", nil)

	user, err := svc.GetByID(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != "u1" {
		t.Fatalf("unexpected user id: %s", user.ID)
	}
}

