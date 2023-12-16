package core

import "context"

type UserRepository interface {
	GetByUserName(ctx context.Context, userName string) (User, error)
	Create(ctx context.Context, user User) error
	Delete(ctx context.Context, userId string) error
}

type TaskRepository interface {
	Create(ctx context.Context, task Task) error
	GetTasksByUserId(ctx context.Context, userId uint) ([]Task, error)
	Edit(ctx context.Context, id uint, update Task) error
	UpdateTaskStatus(ctx context.Context, userId uint, taskId uint) (bool, error)
	SearchTasks(ctx context.Context, priority string, status string, sortType string, searchString string) []Task
}

type UserService interface {
	Create(ctx context.Context, user User) error
	CheckIfUserIsAdmin(ctx context.Context, userName string) (bool, error)
	Delete(ctx context.Context, userId string) error
	GetUserProfileData(ctx context.Context, userName string) (User, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, task Task) error
	UpdateTask(ctx context.Context, taskId uint, task Task) error
	MarkTaskStatus(ctx context.Context, taskId uint, id uint) (bool, error)
	SearchTasks(ctx context.Context, priority string, status string, sortType string, searchString string) []Task
}
