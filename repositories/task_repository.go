package repositories

import (
	"gotask-api/models"

	"gorm.io/gorm"
)

// interface - define what operation are available

type TaskRepository interface {
	FindAll() ([]models.Task, error)
	FindByID(id uint) (*models.Task, error)
	Create(task *models.Task) (*models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(id uint) error
}

// implementasi dari si interfacenya ke DB
type taskRepository struct {
	db *gorm.DB
}

// buat constructornya
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Order("id asc").Find(&tasks)
	return tasks, result.Error
}

func (r *taskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (r *taskRepository) Create(task *models.Task) (*models.Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}
func (r *taskRepository) Update(task *models.Task) (*models.Task, error) {
	result := r.db.Save(task)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}

func (r *taskRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
