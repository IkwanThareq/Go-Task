package domains

import (
	"errors"
	"strings"

	"gotask-api/datatransfers"
	"gotask-api/models"
	"gotask-api/repositories"
)

// buat sentinel error
var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrTaskAlreadyDone = errors.New("task already complete")
	ErrNotTaskOwner    = errors.New("you do not own this task")
)

// interface
type TaskDomain interface {
	CreateTask(userID uint, title string, description string, priority int) (*models.Task, error)
	GetAllTasks(userID uint, params datatransfers.TaskQueryParams) (*datatransfers.PaginatedResponse, error)
	GetTaskById(userID uint, id uint) (*models.Task, error)
	UpdateTask(userID uint, id uint, title string, description string, priority int, status string) (*models.Task, error)
	DeleteTask(userID uint, id uint) error
}

// implementation from the interface

type taskDomain struct {
	repo repositories.TaskRepository
}

// constructor
func NewTaskDomain(repo repositories.TaskRepository) TaskDomain {
	return &taskDomain{repo: repo}
}

func (t *taskDomain) GetAllTasks(userID uint, params datatransfers.TaskQueryParams) (*datatransfers.PaginatedResponse, error) {
	// set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	tasks, totalItems, err := t.repo.FindAllByUserID(
		userID,
		params.Status,
		params.Priority,
		params.Limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	// calculate total pages
	totalPages := int(totalItems) / params.Limit
	if int(totalItems)%params.Limit != 0 {
		totalPages++
	}

	return &datatransfers.PaginatedResponse{
		Items:      tasks,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (t *taskDomain) GetTaskById(userID uint, id uint) (*models.Task, error) {
	task, err := t.repo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, ErrNotTaskOwner
	}
	return task, nil
}

func (t *taskDomain) CreateTask(userID uint, title, description string, priority int) (*models.Task, error) {
	// mulai aturan bisnis, validasi ada disini bukan di handler
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be blank")
	}
	if strings.TrimSpace(description) == "" {
		return nil, errors.New("description cannot be blank")
	}
	if priority < 1 || priority > 3 {
		return nil, errors.New("priority must be between 1 and 3")
	}

	task := &models.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      "pending",
	}

	return t.repo.Create(task)
}

func (d *taskDomain) UpdateTask(userID uint, id uint, title, description string, priority int, status string) (*models.Task, error) {
	// check task exists first
	task, err := d.GetTaskById(userID, id)
	if err != nil {
		return nil, err
	}

	// validate — only update fields that are provided
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if strings.TrimSpace(description) == "" {
		return nil, errors.New("description cannot be empty")
	}
	if priority < 1 || priority > 3 {
		return nil, errors.New("priority must be between 1 and 3")
	}

	// validate status
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"done":        true,
	}
	if !validStatuses[status] {
		return nil, errors.New("status must be pending, in_progress, or done")
	}

	// apply updates
	// task.Title = title
	// task.Description = description
	// task.Priority = priority
	// task.Status = status

	// update fields
	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if priority != 0 {
		task.Priority = priority
	}
	if status != "" {
		task.Status = status
	}

	return d.repo.Update(task)
}

func (t *taskDomain) DeleteTask(userID uint, id uint) error {
	// check exist atau tidak
	_, err := t.GetTaskById(userID, id)
	if err != nil {
		return ErrTaskNotFound
	}

	return t.repo.Delete(id)
}
