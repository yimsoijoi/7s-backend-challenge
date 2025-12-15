package application

import (
	"context"
	"errors"
	"time"

	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository
	jwt  infrastructure.JWTManager
}

func NewUserService(r domain.UserRepository, jwt infrastructure.JWTManager) domain.UserService {
	return &UserService{repo: r, jwt: jwt}
}

func (s *UserService) Register(ctx context.Context, name, email, password string) error {
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

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.Generate(user.ID)
}

func (s *UserService) Update(
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

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
