package repositories

import (
	"context"
	"errors"
	"fmt"
	"taskmanager/core"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (tr TaskRepository) Create(ctx context.Context, task core.Task) error {
	return tr.db.WithContext(ctx).Create(&task).Error
}

func (tr TaskRepository) GetTasksByUserId(ctx context.Context, userId uint) ([]core.Task, error) {
	userTasks := []core.Task{}
	result := tr.db.WithContext(ctx).Where("user_id = ?", userId).Find(&userTasks)
	if result.Error != nil {
		return nil, errors.New("Error while pulling user task")
	}
	return userTasks, nil
}

func (tr TaskRepository) Edit(ctx context.Context, id uint, update core.Task) error {
	var existingTask core.Task
	updateOp := tr.db.WithContext(ctx).Where(&core.Task{
		ID: id,
	}).First(&existingTask)

	if updateOp.Error != nil {
		if errors.Is(updateOp.Error, gorm.ErrRecordNotFound) {
			return errors.New("Task Not Found")
		}
		return updateOp.Error
	}

	if update.UserID != 0 {
		existingTask.UserID = update.UserID
	}
	if update.Title != "" {
		existingTask.Title = update.Title
	}
	if update.Description != "" {
		existingTask.Description = update.Description
	}
	if update.Priority != "" {
		existingTask.Priority = update.Priority
	}
	if update.DueDate != nil {
		existingTask.DueDate = update.DueDate
	}

	return tr.db.WithContext(ctx).Save(&existingTask).Error
}

func (tr TaskRepository) UpdateTaskStatus(ctx context.Context, userId uint, taskId uint) (bool, error) {
	var existingTask core.Task
	result := tr.db.WithContext(ctx).Where(&core.Task{ID: taskId, UserID: userId}).First(&existingTask)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("Task Not Found")
		}
		return false, result.Error
	}

	existingTask.Status = !existingTask.Status
	return existingTask.Status, tr.db.WithContext(ctx).Save(existingTask).Error

}

func (tr TaskRepository) SearchTasks(ctx context.Context, priority string, status string, sortType string, searchString string) []core.Task {
	var tasks []core.Task

	query := tr.db.WithContext(ctx).Model(&core.Task{})
	fmt.Print(searchString, "search string")
	if searchString != "" {
		searchString := "%" + searchString + "%"
		query = query.Where("title like ?", searchString)
		query = query.Where("description like ?", searchString)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if status != "" {
		if status == "true" {
			query = query.Where("status = ?", true)
		} else if status == "false" {
			query = query.Where("status = ?", false)
		}
	}

	if sortType != "" {
		fmt.Println(sortType, "sortType")
		if sortType == "desc" {
			query = query.Order("due_date desc")
		} else {
			query = query.Order("due_date asc")
		}
	}
	fmt.Println(priority)

	query.Find(&tasks)
	return tasks
}
