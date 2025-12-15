package mocks

import (
	"context"
	"errors"

	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
)

type UserRepositoryMock struct {
	FindByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	FindByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	FindAllFn     func(ctx context.Context) ([]*domain.User, error)
	CreateFn      func(ctx context.Context, user *domain.User) error
	UpdateFn      func(ctx context.Context, user *domain.User) error
	DeleteFn      func(ctx context.Context, id string) error
	CountFn       func(ctx context.Context) (int64, error)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *UserRepositoryMock) FindAll(ctx context.Context) ([]*domain.User, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *UserRepositoryMock) Count(ctx context.Context) (int64, error) {
	if m.CountFn != nil {
		return m.CountFn(ctx)
	}
	return 0, errors.New("not implemented")
}
