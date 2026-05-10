package repository

import (
	"context"

	"github.com/zanwyyy/platform/internal/domain/entity"
)

// UserRepository defines the contract for user persistence operations.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
