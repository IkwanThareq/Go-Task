package domains

import (
	"errors"
	"gotask-api/models"
	"gotask-api/repositories"
	"strings"
)

// buat sentinel error
var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrTaskAlreadyDone = errors.New("task already complete")
)

// interface
type TaskDomain interface {
	GetAllTasks() ([]models.Task, error)
	GetTaskById(id uint) (*models.Task, error)
	CreateTask(title, description string, priority int) (*models.Task, error)
	UpdateTask(id uint, title, description string, priority int, status string) (*models.Task, error)
	DeleteTask(id uint) error
}

// implementation from the interface

type taskDomain struct {
	repo repositories.TaskRepository
}

// constructor
func NewTaskDomain(repo repositories.TaskRepository) TaskDomain {
	return &taskDomain{repo: repo}
}

func (t *taskDomain) GetAllTasks() ([]models.Task, error) {
	return t.repo.FindAll()
}

func (t *taskDomain) GetTaskById(id uint) (*models.Task, error) {
	task, err := t.repo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (t *taskDomain) CreateTask(title, description string, priority int) (*models.Task, error) {
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
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      "pending",
	}

	return t.repo.Create(task)
}

func (d *taskDomain) UpdateTask(id uint, title, description string, priority int, status string) (*models.Task, error) {
	// check task exists first
	task, err := d.repo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
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
	task.Title = title
	task.Description = description
	task.Priority = priority
	task.Status = status

	return d.repo.Update(task)
}

func (t *taskDomain) DeleteTask(id uint) error {
	// check exist atau tidak
	_, err := t.repo.FindByID(id)
	if err != nil {
		return ErrTaskNotFound
	}

	return t.repo.Delete(id)
}
