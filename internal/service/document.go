package service

import (
	"fmt"
	"time"
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type DocumentService struct {
	repo repository.DocumentRepositoryInterface
}

func NewDocumentService(repo repository.DocumentRepositoryInterface) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) CreateDocument(doc *model.Document) *model.DocumentResponse {
	// Валидация документа
	if err := s.validateDocument(doc); err != nil {
		return &model.DocumentResponse{
			Message: err.Error(),
		}
	}

	// Генерируем номер документа
	doc.Number = s.generateDocumentNumber(doc.Type)

	// Устанавливаем статус черновика
	doc.Status = "draft"

	// Если дата не установлена, используем текущую
	if doc.Date.Time().IsZero() {
		doc.Date = model.Date(time.Now())
	}

	// Рассчитываем TotalPrice для каждой позиции
	for i := range doc.Items {
		doc.Items[i].TotalPrice = doc.Items[i].Price * float64(doc.Items[i].Quantity)
	}

	response := s.repo.CreateDocument(doc)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) GetDocument(id int) *model.DocumentResponse {
	response := s.repo.GetDocument(id)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) GetAllDocuments() *model.DocumentListResponse {
	response := s.repo.GetAllDocuments()
	return &model.DocumentListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) UpdateDocument(doc *model.Document) *model.DocumentResponse {
	response := s.repo.UpdateDocument(doc)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) DeleteDocument(id int) *model.DocumentResponse {
	response := s.repo.DeleteDocument(id)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) ApproveDocument(id int, approvedByID int) *model.DocumentResponse {
	response := s.repo.ApproveDocument(id, approvedByID)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

// Вспомогательные методы

func (s *DocumentService) validateDocument(doc *model.Document) error {
	if doc.Type == "" {
		return fmt.Errorf("тип документа не указан")
	}

	if doc.LocationID == 0 {
		return fmt.Errorf("местоположение не указано")
	}

	if doc.CreatedByID == 0 {
		return fmt.Errorf("создатель документа не указан")
	}

	if len(doc.Items) == 0 {
		return fmt.Errorf("документ должен содержать хотя бы одну позицию")
	}

	// Проверяем позиции документа
	for i, item := range doc.Items {
		if item.EquipmentID == 0 {
			return fmt.Errorf("оборудование не указано в позиции %d", i+1)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("количество должно быть больше нуля в позиции %d", i+1)
		}

		if doc.Type == "inventory" && item.ActualQuantity < 0 {
			return fmt.Errorf("фактическое количество не может быть отрицательным в позиции %d", i+1)
		}

		if item.Price < 0 {
			return fmt.Errorf("цена не может быть отрицательной в позиции %d", i+1)
		}
	}

	return nil
}

func (s *DocumentService) generateDocumentNumber(docType string) string {
	// Получаем префикс в зависимости от типа документа
	prefix := map[string]string{
		"inventory":  "ИНВ",
		"transfer":   "ПЕР",
		"write_off":  "СПС",
		"acceptance": "ПРМ",
	}[docType]

	// Получаем текущий год
	year := time.Now().Year()

	// Получаем все документы этого типа
	response := s.GetAllDocuments()
	count := 1
	if response.Model != nil {
		// Считаем количество документов этого типа
		for _, doc := range response.Model {
			if doc.Type == docType {
				count++
			}
		}
	}

	// Формируем номер в формате ПРФ-2024-001
	return fmt.Sprintf("%s-%d-%03d", prefix, year, count)
}
