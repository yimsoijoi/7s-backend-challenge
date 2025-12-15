package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yimsoijoi/7s-backend-challenge/internal/application"
	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"golang.org/x/crypto/bcrypt"

	"github.com/yimsoijoi/7s-backend-challenge/internal/adapters/mongo/mocks"
	jwtmocks "github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure/mocks"
)

func TestUserService_Register_Success(t *testing.T) {

	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, "John", user.Name)
			assert.Equal(t, "john@test.com", user.Email)
			assert.NotEmpty(t, user.Password)
			return nil
		},
	}

	jwt := &jwtmocks.JWTManagerMock{}

	svc := application.NewUserService(repo, jwt)

	err := svc.Register(context.Background(), "John", "john@test.com", "secret")

	assert.NoError(t, err)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{}, nil
		},
	}

	svc := application.NewUserService(repo, &jwtmocks.JWTManagerMock{})

	err := svc.Register(context.Background(), "John", "john@test.com", "secret")

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
}

func TestUserService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:       "user-id",
				Password: string(hash),
			}, nil
		},
	}

	jwt := &jwtmocks.JWTManagerMock{
		GenerateFn: func(userID string) (string, error) {
			assert.Equal(t, "user-id", userID)
			return "jwt-token", nil
		},
	}

	svc := application.NewUserService(repo, jwt)

	token, err := svc.Login(context.Background(), "john@test.com", "secret")

	assert.NoError(t, err)
	assert.Equal(t, "jwt-token", token)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	repo := &mocks.UserRepositoryMock{
		FindByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{Password: string(hash)}, nil
		},
	}

	svc := application.NewUserService(repo, &jwtmocks.JWTManagerMock{})

	_, err := svc.Login(context.Background(), "john@test.com", "wrong")

	assert.Error(t, err)
}

func TestUserService_Update(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		FindByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id}, nil
		},
		UpdateFn: func(ctx context.Context, user *domain.User) error {
			assert.Equal(t, "New", user.Name)
			assert.Equal(t, "new@test.com", user.Email)
			return nil
		},
	}

	svc := application.NewUserService(repo, &jwtmocks.JWTManagerMock{})

	err := svc.Update(context.Background(), "id", "New", "new@test.com")

	assert.NoError(t, err)
}
