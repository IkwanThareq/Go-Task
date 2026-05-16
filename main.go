package main

import (
	"log"
	"net/http"
	"os"

	"gotask-api/config"
	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/domains"
	"gotask-api/handlers"
	"gotask-api/models"
	"gotask-api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	// connect DB
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Task{})

	// wire layers together — dependency injection
	taskRepo := repositories.NewTaskRepository(config.DB)
	taskDomain := domains.NewTaskDomain(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskDomain)

	// router
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, gin.H{
				"status":  "ok",
				"service": "gotask-api",
			})
		})

		tasks := v1.Group("/tasks")
		{
			tasks.GET("", taskHandler.HandleGetAllTasks)
			tasks.POST("", taskHandler.HandleCreateTask)
			tasks.GET("/:id", taskHandler.HandleGetTaskByID)
			tasks.PUT("/:id", taskHandler.HandleUpdateTask)
			tasks.DELETE("/:id", taskHandler.HandleDeleteTask)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		datatransfers.ErrorRes(c, http.StatusNotFound,
			constants.NOT_FOUND, "route not found")
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
