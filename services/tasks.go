package services

import (
	"context"
	"taskmanager/core"
)

type TaskService struct {
	TaskRepo core.TaskRepository
}

func NewTaskService(t core.TaskRepository) *TaskService {
	return &TaskService{
		TaskRepo: t,
	}
}

func (t TaskService) CreateTask(ctx context.Context, task core.Task) error {
	return t.TaskRepo.Create(ctx, task)
}

func (t TaskService) UpdateTask(ctx context.Context, taskId uint, task core.Task) error {
	return t.TaskRepo.Edit(ctx, taskId, task)
}

func (t TaskService) MarkTaskStatus(ctx context.Context, taskId uint, id uint) (bool, error) {
	isUpdated, err := t.TaskRepo.UpdateTaskStatus(ctx, id, taskId)
	return isUpdated, err
}

func (t TaskService) SearchTasks(ctx context.Context, priority string, statusStr string, sortType string, searchString string) []core.Task {
	return t.TaskRepo.SearchTasks(ctx, priority, statusStr, sortType, searchString)
}
