package memory

import (
	"context"
	"sync"

	"github.com/zanwyyy/platform/internal/domain/entity"
)

// UserRepository is a thread-safe, in-memory implementation of domain/repository.UserRepository.
type UserRepository struct {
	mu    sync.RWMutex
	store map[string]*entity.User
}

// NewUserRepository returns an empty in-memory UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{store: make(map[string]*entity.User)}
}

func (r *UserRepository) FindByID(_ context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.store[id], nil
}

func (r *UserRepository) FindAll(_ context.Context) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]*entity.User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Save(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = user
	return nil
}

func (r *UserRepository) Update(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = user
	return nil
}

func (r *UserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}
