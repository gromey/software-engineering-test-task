package handler

import (
	"cruder/internal/controller"
	"cruder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(router *gin.Engine, apiKey string, userController *controller.UserController) *gin.Engine {
	v1 := router.Group("/api/v1", middleware.APIKey(apiKey), middleware.Logging)
	{
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/", userController.GetAllUsers)
			userGroup.GET("/username/:username", userController.GetUserByUsername)
			userGroup.GET("/id/:id", userController.GetUserByID)
			userGroup.POST("/", userController.PostUser)
			userGroup.PATCH("/:id", userController.PatchUser)
			userGroup.DELETE("/:id", userController.DeleteUser)
		}
	}
	return router
}
