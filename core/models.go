package core

import (
	"time"

	"gorm.io/gorm"
)

// User hold FK Task model

type Task struct {
	gorm.Model
	ID          uint `json:"id" gorm:"primary_key; unique"`
	UserID      uint
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
	Status      bool       `json:"status" gorm:"default:false"`
}

type User struct {
	gorm.Model
	Username string `json:"userName" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Type     string `gorm:"default:'default'"`
	Tasks    []Task `gorm:"foreignKey:UserID"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Type     string `json:"type"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type CreateTaskRequest struct {
	UserID      uint       `json:"userId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
	Status      bool       `json:"status"`
}
