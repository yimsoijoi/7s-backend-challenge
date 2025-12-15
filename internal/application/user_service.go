package application

import (
	"context"
	"errors"
	"time"

	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure"
	"github.com/yimsoijoi/7s-backend-challenge/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo ports.UserRepository
	jwt  infrastructure.JWTManager
}

func NewUserService(r ports.UserRepository, jwt infrastructure.JWTManager) ports.UserService {
	return &userService{repo: r, jwt: jwt}
}

func (s *userService) Register(ctx context.Context, name, email, password string) error {
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return errors.New("email already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		Name:      name,
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.Generate(user.ID)
}

func (s *userService) Update(
	ctx context.Context,
	id, name, email string,
) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	user.Name = name
	user.Email = email
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
