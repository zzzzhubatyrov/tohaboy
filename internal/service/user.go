package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(username string) model.Response[*model.User] {
	return s.repo.GetUser(username)
}

func (s *UserService) GetCurrentUser() model.Response[*model.User] {
	// В данной реализации возвращаем админа, так как у нас пока нет сессий
	return s.repo.GetUser("admin")
}

func (s *UserService) GetByID(id int) model.Response[*model.User] {
	return s.repo.GetByID(id)
}

func (s *UserService) GetByName(name string) model.Response[*model.User] {
	return s.repo.GetUser(name)
}

func (s *UserService) Update(user *model.User) model.Response[*model.User] {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id int) model.Response[*model.User] {
	return s.repo.Delete(id)
}

// Вспомогательные функции

func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"admin":   true,
		"manager": true,
		"auditor": true,
	}
	return validRoles[role]
}
