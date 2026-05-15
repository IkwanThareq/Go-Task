package main

import (
	"gotask-api/config"
	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var tasks = []gin.H{}
var lastID = 0

func main() {
	// load .env file data
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// connect to DB
	config.ConnectDatabase()

	// auto migrate - crates tables from struct
	config.DB.AutoMigrate(&models.Task{})

	// create gin router
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
			tasks.GET("", handleGetAllTasks)
			tasks.POST("", handleCreateTask)
			tasks.GET("/:id", handleGetTaskByID)
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

// all the function for the endpoint

func handleGetAllTasks(c *gin.Context) {
	var tasks []models.Task
	result := config.DB.Find(&tasks)
	if result.Error != nil {
		datatransfers.ErrorRes(c, http.StatusInternalServerError,
			constants.INTERNAL_ERROR, "failed to fetch tasks")
		return
	}
	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, tasks)
}

func handleCreateTask(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "invalid request body")
		return
	}

	if body.Title == "" {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "title cannot be empty")
		return
	}

	// adding description as mandatory
	if strings.TrimSpace(body.Description) == "" {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "description cannot be empty")
		return
	}

	if body.Priority < 1 || body.Priority > 3 {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "priority must be between 1 and 3")
		return
	}

	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		Status:      constants.StatusPending,
	}

	result := config.DB.Create(&task)
	if result.Error != nil {
		datatransfers.ErrorRes(c, http.StatusInternalServerError,
			constants.INTERNAL_ERROR, "failed to create task")
		return
	}

	datatransfers.SuccessRes(c, http.StatusCreated, constants.SUCCESS, task)
}

func handleGetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "invalid task id")
		return
	}

	var task models.Task
	result := config.DB.First(&task, id)
	if result.Error != nil {
		datatransfers.ErrorRes(c, http.StatusNotFound,
			constants.NOT_FOUND, "task not found")
		return
	}

	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, task)
}
