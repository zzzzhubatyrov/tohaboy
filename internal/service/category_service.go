package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepositoryInterface
}

func NewCategoryService(repo repository.CategoryRepositoryInterface) CategoryServiceInterface {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(category *model.Category) *model.CategoryResponse {
	response := s.repo.CreateCategory(category)
	return &model.CategoryResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *CategoryService) GetCategory(id int) *model.CategoryResponse {
	response := s.repo.GetCategory(id)
	return &model.CategoryResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *CategoryService) GetAllCategories() *model.CategoryListResponse {
	response := s.repo.GetAllCategories()
	return &model.CategoryListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *CategoryService) UpdateCategory(category *model.Category) *model.CategoryResponse {
	response := s.repo.UpdateCategory(category)
	return &model.CategoryResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *CategoryService) DeleteCategory(id int) *model.CategoryResponse {
	response := s.repo.DeleteCategory(id)
	return &model.CategoryResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}
