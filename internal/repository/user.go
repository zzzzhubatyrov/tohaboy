package repository

import (
	"errors"
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(username string) model.Response[*model.User] {
	var user model.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Response[*model.User]{
				Message: "Пользователь не найден",
			}
		}
		return model.Response[*model.User]{
			Message: "Ошибка при получении пользователя",
		}
	}

	return model.Response[*model.User]{
		Message: "Пользователь успешно получен",
		Model:   &user,
	}
}

func (r *UserRepository) GetByID(id int) model.Response[*model.User] {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Response[*model.User]{
				Message: "Пользователь не найден",
			}
		}
		return model.Response[*model.User]{
			Message: "Ошибка при получении пользователя",
		}
	}

	return model.Response[*model.User]{
		Message: "Пользователь успешно получен",
		Model:   &user,
	}
}

func (r *UserRepository) Update(user *model.User) model.Response[*model.User] {
	result := r.db.Save(user)
	if result.Error != nil {
		return model.Response[*model.User]{
			Message: "Ошибка при обновлении пользователя",
		}
	}

	return model.Response[*model.User]{
		Message: "Пользователь успешно обновлен",
		Model:   user,
	}
}

func (r *UserRepository) Delete(id int) model.Response[*model.User] {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Response[*model.User]{
				Message: "Пользователь не найден",
			}
		}
		return model.Response[*model.User]{
			Message: "Ошибка при получении пользователя",
		}
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return model.Response[*model.User]{
			Message: "Ошибка при удалении пользователя",
		}
	}

	return model.Response[*model.User]{
		Message: "Пользователь успешно удален",
		Model:   &user,
	}
}
