package usecase_test

import (
	"context"
	"testing"

	"github.com/zanwyyy/platform/internal/repository/memory"
	"github.com/zanwyyy/platform/internal/usecase"
)

func newUserUC() usecase.UserUseCase {
	return usecase.NewUserUseCase(memory.NewUserRepository())
}

func TestCreate(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	user, err := uc.Create(ctx, usecase.CreateUserInput{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == "" {
		t.Error("expected non-empty ID")
	}
	if user.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", user.Email)
	}
}

func TestGetByID(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	created, _ := uc.Create(ctx, usecase.CreateUserInput{Name: "Bob", Email: "bob@example.com"})

	got, err := uc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, got.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	_, err := uc.GetByID(ctx, "non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != usecase.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetAll(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	uc.Create(ctx, usecase.CreateUserInput{Name: "A", Email: "a@example.com"}) //nolint:errcheck
	uc.Create(ctx, usecase.CreateUserInput{Name: "B", Email: "b@example.com"}) //nolint:errcheck

	users, err := uc.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUpdate(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	created, _ := uc.Create(ctx, usecase.CreateUserInput{Name: "Carol", Email: "carol@example.com"})

	updated, err := uc.Update(ctx, usecase.UpdateUserInput{
		ID:    created.ID,
		Name:  "Carol Updated",
		Email: "carol2@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Carol Updated" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	_, err := uc.Update(ctx, usecase.UpdateUserInput{ID: "ghost", Name: "X", Email: "x@example.com"})
	if err != usecase.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	created, _ := uc.Create(ctx, usecase.CreateUserInput{Name: "Dave", Email: "dave@example.com"})

	if err := uc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := uc.GetByID(ctx, created.ID)
	if err != usecase.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	uc := newUserUC()
	ctx := context.Background()

	err := uc.Delete(ctx, "non-existent")
	if err != usecase.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
