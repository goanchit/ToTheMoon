package services

import (
	"context"
	"errors"
	"fmt"
	"taskmanager/core"
)

type UserService struct {
	UserRepo core.UserRepository
	TaskRepo core.TaskRepository
}

func NewUserService(u core.UserRepository, t core.TaskRepository) *UserService {
	return &UserService{
		UserRepo: u,
		TaskRepo: t,
	}
}

func (u UserService) Create(ctx context.Context, user core.User) error {
	return u.UserRepo.Create(ctx, user)
}

func (u UserService) CheckIfUserIsAdmin(ctx context.Context, userName string) (bool, error) {
	userInfo, err := u.UserRepo.GetByUserName(ctx, userName)
	if err != nil {
		return false, errors.New(err.Error())
	}
	if userInfo.Type == "admin" {
		return true, nil
	}
	return false, nil
}

func (u UserService) Delete(ctx context.Context, userId string) error {
	return u.UserRepo.Delete(ctx, userId)
}

func (u UserService) GetUserProfileData(ctx context.Context, userName string) (core.User, error) {
	userInfo, err := u.UserRepo.GetByUserName(ctx, userName)
	if err != nil {
		return core.User{}, errors.New(err.Error())
	}
	// Get user related tasks
	userTasks, err := u.TaskRepo.GetTasksByUserId(ctx, userInfo.ID)
	fmt.Print(userTasks, "userTasksuserTasksuserTasks")
	if err != nil {
		return userInfo, err
	}
	userInfo.Tasks = userTasks
	return userInfo, nil
}
