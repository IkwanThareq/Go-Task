package main

import (
	"gotask-api/config"
	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/models"
	"log"
	"net/http"
	"strconv"

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

	//// define route
	//// GET /ping -> runs the function below
	//router.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})
	//
	//// health check
	//router.GET("/health", func(c *gin.Context) {
	//	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, gin.H{
	//		"status":  "ok",
	//		"service": "gotask-api",
	//	})
	//})
	//
	//router.NoRoute(func(c *gin.Context) {
	//	datatransfers.ErrorRes(c, http.StatusNotFound, constants.NOT_FOUND, "route not found")
	//})
	//
	//// creating Route Group for the API now is using v1
	//v1 := router.Group("/api/v1"){
	//	tasks := v1.Group("/tasks"){
	//		tasks.GET("", getAllTasks)
	//		tasks.POST("", createTasks)
	//		tasks.GET("/:id", getTaskByID)
	//		tasks.PUT("/:id", updateTask)
	//		tasks.DELETE("/:id", deleteTask)
	//	}
	//}

	// above is for tranining purpose
	// and this below is for the course training
	// in memory storage

	v1 := router.Group("/api/v1")
	{
		// health check
		v1.GET("/health", func(c *gin.Context) {
			datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, gin.H{
				"status": "ok",
				"server": "gotask-api",
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
		datatransfers.ErrorRes(c, http.StatusNotFound, constants.NOT_FOUND, "route not found")
	})

	router.Run(":8080")
}

// all the function for the endpoint

func handleGetAllTasks(c *gin.Context) {
	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, tasks)
}

func handleCreateTask(c *gin.Context) {
	var body struct {
		Title    string `json:"title"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, "invalid request body")
		return
	}

	if body.Title == "" {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, "title is required")
		return
	}

	// start logic fot create task
	lastID++
	newTask := gin.H{
		"id":       lastID,
		"title":    body.Title,
		"priority": body.Priority,
		"status":   constants.StatusPending,
	}
	tasks = append(tasks, newTask)
	datatransfers.SuccessRes(c, http.StatusCreated, constants.SUCCESS, newTask)
}

func handleGetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, "invalid task id")
		return
	}

	for _, task := range tasks {
		if task["id"] == id {
			datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, task)
			return
		}
	}

	datatransfers.ErrorRes(c, http.StatusNotFound, constants.NOT_FOUND, "task not found")

}
