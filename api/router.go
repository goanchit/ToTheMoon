package api

import (
	"taskmanager/repositories"
	"taskmanager/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteHandler(g *gin.Engine, db *gorm.DB) {
	api := g.Group("/api")

	userRepo := repositories.NewUserRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	userService := services.NewUserService(userRepo, taskRepo)
	taskService := services.NewTaskService(taskRepo)

	handler := NewAPIHandler(userService, taskService)

	userHandler := api.Group("/user")
	{
		userHandler.POST("/", handler.CreateUser)
		userHandler.GET("/", handler.GetProfileData)
		userHandler.DELETE("/", handler.DeleteUser)
	}

	taskHandler := api.Group("/task")
	{
		taskHandler.POST("/", handler.CreateTask)
		taskHandler.PATCH("/:id", handler.UpdateTask)
		taskHandler.PATCH("/updateTask/:id", handler.MarkTaskStatus)
		taskHandler.GET("/search", handler.SearchTask)
	}

}
