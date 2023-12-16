package repositories

import (
	"context"
	"errors"
	"taskmanager/core"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur UserRepository) GetByUserName(ctx context.Context, userName string) (core.User, error) {
	user := core.User{}
	result := ur.db.WithContext(ctx).Where(core.User{
		Username: userName,
	}).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("User not found")
	}
	return user, result.Error

}

func (ur UserRepository) Create(ctx context.Context, usr core.User) error {
	return ur.db.WithContext(ctx).Create(&usr).Error
}

func (ur UserRepository) Delete(ctx context.Context, userId string) error {
	return ur.db.WithContext(ctx).Delete(&core.User{}, userId).Error
}
