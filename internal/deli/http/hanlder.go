package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zanwyyy/platform/internal/builder"
)

type Handler struct {
	builderSvc *builder.Service
}

func NewHandler(builderSvc *builder.Service) *Handler {
	return &Handler{builderSvc: builderSvc}
}

func (h *Handler) CreateBuild(c *gin.Context) {
	var req builder.BuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chạy build dưới dạng Background Job để không block API
	go func() {
		err := h.builderSvc.ProcessBuild(c, req)
		if err != nil {
			// Ở đây bạn có thể gọi Webhook ngược lại Platform API để báo lỗi
			println("Build Error:", err.Error())
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Build task started in background",
		"status":  "processing",
	})
}
