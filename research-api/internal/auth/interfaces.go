package auth

import "context"

type Service interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, userID string) (*User, error)
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, userID string) (*User, error)
}
