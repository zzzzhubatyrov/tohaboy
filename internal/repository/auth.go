package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) Register(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepo) Login(user *model.User) (*model.User, error) {
	var existingUser model.User
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return nil, err
	}
	return &existingUser, nil
}
