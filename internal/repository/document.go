package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func (r *DocumentRepository) GetDB() *gorm.DB {
	return r.db
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) CreateDocument(doc *model.Document) model.Response[*model.Document] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	// Создаем документ без items
	items := doc.Items
	doc.Items = nil
	if err := tx.Create(doc).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Создаем позиции документа
	for i := range items {
		items[i].DocumentID = doc.ID
		items[i].ID = 0 // Ensure ID is zero to let GORM auto-increment
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return model.Response[*model.Document]{
				Message: err.Error(),
			}
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Загружаем созданный документ со всеми связями
	if err := r.db.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		First(doc, doc.ID).Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Document]{
		Model: doc,
	}
}

func (r *DocumentRepository) GetDocument(id uint) model.Response[*model.Document] {
	var doc model.Document

	if err := r.db.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		First(&doc, id).Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Document]{
		Model: &doc,
	}
}

func (r *DocumentRepository) UpdateDocument(doc *model.Document) model.Response[*model.Document] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	// Проверяем существование документа
	var existingDoc model.Document
	if err := tx.First(&existingDoc, doc.ID).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: "Документ не найден",
		}
	}

	// Проверяем статус документа
	if existingDoc.Status != "draft" {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: "Можно редактировать только черновики",
		}
	}

	// Удаляем старые позиции
	if err := tx.Where("document_id = ?", doc.ID).Delete(&model.DocumentItem{}).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Обновляем документ
	if err := tx.Save(doc).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Создаем новые позиции
	for i := range doc.Items {
		doc.Items[i].DocumentID = doc.ID
		doc.Items[i].ID = 0 // Reset ID to let GORM auto-increment
		if err := tx.Create(&doc.Items[i]).Error; err != nil {
			tx.Rollback()
			return model.Response[*model.Document]{
				Message: err.Error(),
			}
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Загружаем обновленный документ со всеми связями
	if err := r.db.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		First(doc, doc.ID).Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Document]{
		Model: doc,
	}
}

func (r *DocumentRepository) DeleteDocument(id uint) model.Response[*model.Document] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	// Получаем документ для проверки статуса и возврата данных
	var doc model.Document
	if err := tx.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		First(&doc, id).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Проверяем статус документа
	if doc.Status != "draft" {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: "Можно удалять только черновики",
		}
	}

	// Удаляем позиции документа
	if err := tx.Where("document_id = ?", id).Delete(&model.DocumentItem{}).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Удаляем сам документ
	if err := tx.Delete(&doc).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Document]{
		Model: &doc,
	}
}

func (r *DocumentRepository) GetAllDocuments() model.Response[[]model.Document] {
	var docs []model.Document

	if err := r.db.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Find(&docs).Error; err != nil {
		return model.Response[[]model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Document]{
		Model: docs,
	}
}

func (r *DocumentRepository) ApproveDocument(id uint, approvedByID uint) model.Response[*model.Document] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	// Получаем документ
	var doc model.Document
	if err := tx.First(&doc, id).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Проверяем статус документа
	if doc.Status != "draft" {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: "Можно утвердить только черновик",
		}
	}

	// Обновляем статус и добавляем утверждающего
	doc.Status = "completed"
	doc.ApprovedByID = approvedByID

	if err := tx.Save(&doc).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Загружаем обновленный документ со всеми связями
	if err := r.db.Preload("Items.Equipment").
		Preload("Location").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		First(&doc, id).Error; err != nil {
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Document]{
		Model: &doc,
	}
}
