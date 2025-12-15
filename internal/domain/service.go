package domain

import "context"

type UserService interface {
	Register(ctx context.Context, name, email, password string) error
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Login(ctx context.Context, email, password string) (string, error)
	Update(ctx context.Context, id, name, email string) error
	Delete(ctx context.Context, id string) error
}
