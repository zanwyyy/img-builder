package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/zanwyyy/platform/internal/domain/entity"
	"github.com/zanwyyy/platform/internal/domain/repository"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// CreateUserInput holds the data required to create a new user.
type CreateUserInput struct {
	Name  string
	Email string
}

// UpdateUserInput holds the data required to update an existing user.
type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
}

// UserUseCase defines the business operations available for users.
type UserUseCase interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	Create(ctx context.Context, input CreateUserInput) (*entity.User, error)
	Update(ctx context.Context, input UpdateUserInput) (*entity.User, error)
	Delete(ctx context.Context, id string) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase backed by the given repository.
func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userUseCase) GetAll(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *userUseCase) Create(ctx context.Context, input CreateUserInput) (*entity.User, error) {
	now := time.Now().UTC()
	user := &entity.User{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) Update(ctx context.Context, input UpdateUserInput) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Name = input.Name
	user.Email = input.Email
	user.UpdatedAt = time.Now().UTC()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) Delete(ctx context.Context, id string) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return u.userRepo.Delete(ctx, id)
}
