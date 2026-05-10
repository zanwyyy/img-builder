package router

import (
	"github.com/gin-gonic/gin"

	"github.com/zanwyyy/platform/internal/delivery/http/handler"
)

// Setup registers all application routes onto the given Gin engine.
func Setup(engine *gin.Engine, userHandler *handler.UserHandler) {
	api := engine.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}
}
