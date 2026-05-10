package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zanwyyy/platform/internal/usecase"
	"github.com/zanwyyy/platform/pkg/response"
)

// UserHandler holds the HTTP handlers for user-related endpoints.
type UserHandler struct {
	userUC usecase.UserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userUC usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// createUserRequest is the expected request body for creating a user.
type createUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// updateUserRequest is the expected request body for updating a user.
type updateUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// GetAll godoc
// @Summary     List all users
// @Description Returns every user stored in the system.
// @Tags        users
// @Produce     json
// @Success     200 {object} map[string]any
// @Router      /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userUC.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

// GetByID godoc
// @Summary     Get a user by ID
// @Description Returns the user with the specified ID.
// @Tags        users
// @Produce     json
// @Param       id  path     string true "User ID"
// @Success     200 {object} map[string]any
// @Failure     404 {object} map[string]any
// @Router      /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userUC.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, user)
}

// Create godoc
// @Summary     Create a user
// @Description Stores a new user and returns the created resource.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       body body     createUserRequest true "User payload"
// @Success     201  {object} map[string]any
// @Failure     400  {object} map[string]any
// @Router      /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.userUC.Create(c.Request.Context(), usecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(c, user)
}

// Update godoc
// @Summary     Update a user
// @Description Replaces the name and email of an existing user.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path     string            true "User ID"
// @Param       body body     updateUserRequest true "User payload"
// @Success     200  {object} map[string]any
// @Failure     400  {object} map[string]any
// @Failure     404  {object} map[string]any
// @Router      /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.userUC.Update(c.Request.Context(), usecase.UpdateUserInput{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, user)
}

// Delete godoc
// @Summary     Delete a user
// @Description Removes the user with the specified ID.
// @Tags        users
// @Produce     json
// @Param       id  path string true "User ID"
// @Success     204
// @Failure     404 {object} map[string]any
// @Router      /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUC.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.NoContent(c)
}
