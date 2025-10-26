package service

import (
	"context"

	"cruder/internal/model"
	"cruder/internal/repository"
	"cruder/pkg/validation"
)

type UserService interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Post(ctx context.Context, user *model.User) (int64, error)
	Patch(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if err := validation.ValidateUsername(username); err != nil {
		return nil, err
	}
	return s.repo.GetByUsername(ctx, username)
}

func (s *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) Post(ctx context.Context, user *model.User) (int64, error) {
	if err := validation.ValidateUser(user); err != nil {
		return 0, err
	}
	return s.repo.Post(ctx, user)
}

func (s *userService) Patch(ctx context.Context, user *model.User) error {
	if err := validation.ValidateID(user.ID); err != nil {
		return err
	}
	if err := validation.ValidateUser(user); err != nil {
		return err
	}
	return s.repo.Patch(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if err := validation.ValidateID(id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
