package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/zanwyyy/platform/internal/delivery/http/handler"
	"github.com/zanwyyy/platform/internal/delivery/http/router"
	memrepo "github.com/zanwyyy/platform/internal/repository/memory"
	"github.com/zanwyyy/platform/internal/usecase"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	userRepo := memrepo.NewUserRepository()
	userUC := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUC)
	router.Setup(engine, userHandler)
	return engine
}

func TestGetAll_Empty(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCreate_Success(t *testing.T) {
	r := setupRouter()

	body, _ := json.Marshal(map[string]string{"name": "Alice", "email": "alice@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreate_InvalidBody(t *testing.T) {
	r := setupRouter()

	body, _ := json.Marshal(map[string]string{"name": "NoEmail"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/non-existent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string]string{"name": "X", "email": "x@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/ghost", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDelete_NotFound(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/ghost", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestFullCRUD(t *testing.T) {
	r := setupRouter()

	// Create
	body, _ := json.Marshal(map[string]string{"name": "Bob", "email": "bob@example.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", w.Code)
	}

	var createResp map[string]any
	json.Unmarshal(w.Body.Bytes(), &createResp) //nolint:errcheck
	data := createResp["data"].(map[string]any)
	id := data["id"].(string)

	// Get by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("get: expected 200, got %d", w.Code)
	}

	// Update
	body, _ = json.Marshal(map[string]string{"name": "Bob Updated", "email": "bob2@example.com"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("update: expected 200, got %d", w.Code)
	}

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", w.Code)
	}

	// Verify deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("get after delete: expected 404, got %d", w.Code)
	}
}
