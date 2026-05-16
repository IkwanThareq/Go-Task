package handlers

import (
	"errors"
	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/domains"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	domain domains.TaskDomain
}

// constructor

func NewTaskHandler(domain domains.TaskDomain) *TaskHandler {
	return &TaskHandler{domain: domain}
}

func (h *TaskHandler) HandleGetAllTasks(c *gin.Context) {
	tasks, err := h.domain.GetAllTasks()
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusInternalServerError, constants.INTERNAL_ERROR)
		return
	}
	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, tasks)
}

func (h *TaskHandler) HandleGetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, "invalid task id")
		return
	}

	task, err := h.domain.GetTaskById(uint(id))
	if err != nil {
		if errors.Is(err, domains.ErrTaskNotFound) {
			datatransfers.ErrorRes(c, http.StatusNotFound,
				constants.NOT_FOUND, "task not found")
			return
		}
		datatransfers.ErrorRes(c, http.StatusInternalServerError,
			constants.INTERNAL_ERROR, "failed to fetch task")
		return
	}

	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, task)
}

func (h *TaskHandler) HandleCreateTask(c *gin.Context) {
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

	task, err := h.domain.CreateTask(body.Title, body.Description, body.Priority)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, err.Error())
		return
	}

	datatransfers.SuccessRes(c, http.StatusCreated, constants.SUCCESS, task)
}

func (h *TaskHandler) HandleUpdateTask(c *gin.Context) {
	// parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "invalid task id")
		return
	}

	// parse request body
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, "invalid request body")
		return
	}

	// call domain
	task, err := h.domain.UpdateTask(
		uint(id),
		body.Title,
		body.Description,
		body.Priority,
		body.Status,
	)
	if err != nil {
		if errors.Is(err, domains.ErrTaskNotFound) {
			datatransfers.ErrorRes(c, http.StatusNotFound,
				constants.NOT_FOUND, "task not found")
			return
		}
		datatransfers.ErrorRes(c, http.StatusBadRequest,
			constants.BAD_REQUEST, err.Error())
		return
	}

	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, task)
}

func (h *TaskHandler) HandleDeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, "invalid task id")
	}

	err = h.domain.DeleteTask(uint(id))
	if err != nil {
		if errors.Is(err, domains.ErrTaskNotFound) {
			datatransfers.ErrorRes(c, http.StatusNotFound, constants.NOT_FOUND, "task not found")
			return
		}
		datatransfers.ErrorRes(c, http.StatusInternalServerError, constants.INTERNAL_ERROR, "failed to delete task")
		return
	}
	datatransfers.SuccessRes(c, http.StatusOK, constants.SUCCESS, gin.H{"message": "task deleted successfully"})

}
