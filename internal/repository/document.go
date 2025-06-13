package repository

import (
	"gorm.io/gorm"
	"tohaboy/internal/model"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) CreateDocument(doc *model.Document) model.Response[*model.Document] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	if err := tx.Create(doc).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Document]{
			Message: err.Error(),
		}
	}

	// Создаем позиции документа
	for i := range doc.Items {
		doc.Items[i].DocumentID = doc.ID
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

func (r *DocumentRepository) GetDocument(id int) model.Response[*model.Document] {
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

func (r *DocumentRepository) DeleteDocument(id int) model.Response[*model.Document] {
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

func (r *DocumentRepository) ApproveDocument(id int, approvedByID int) model.Response[*model.Document] {
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
