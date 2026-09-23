package repositories

import (
	"gotask-api/models"

	"gorm.io/gorm"
)

// interface - define what operation are available

type TaskRepository interface {
	FindAll() ([]models.Task, error)
	FindByID(id uint) (*models.Task, error)
	FindAllByUserID(userID uint, status string, priority string, limit int, offset int) ([]models.Task, int64, error)
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

func (r *taskRepository) FindAllByUserID(userID uint, status string, priority string, limit int, offset int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var totalItems int64

	// start building the query
	query := r.db.Where("user_id = ?", userID)

	// apply filters if provided

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// count total item before pagination
	query.Model(&models.Task{}).Count(&totalItems)

	// apply and update the code for pagination
	result := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tasks)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return tasks, totalItems, nil
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
