package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *model.Category) model.Response[*model.Category] {
	if err := r.db.Create(category).Error; err != nil {
		return model.Response[*model.Category]{Message: err.Error()}
	}
	return model.Response[*model.Category]{Model: category}
}

func (r *CategoryRepository) GetCategory(id int) model.Response[*model.Category] {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return model.Response[*model.Category]{Message: err.Error()}
	}
	return model.Response[*model.Category]{Model: &category}
}

func (r *CategoryRepository) GetAllCategories() model.Response[[]model.Category] {
	var categories []model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return model.Response[[]model.Category]{Message: err.Error()}
	}
	return model.Response[[]model.Category]{Model: categories}
}

func (r *CategoryRepository) UpdateCategory(category *model.Category) model.Response[*model.Category] {
	if err := r.db.Save(category).Error; err != nil {
		return model.Response[*model.Category]{Message: err.Error()}
	}
	return model.Response[*model.Category]{Model: category}
}

func (r *CategoryRepository) DeleteCategory(id int) model.Response[*model.Category] {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return model.Response[*model.Category]{Message: err.Error()}
	}
	if err := r.db.Delete(&category).Error; err != nil {
		return model.Response[*model.Category]{Message: err.Error()}
	}
	return model.Response[*model.Category]{Model: &category}
}
