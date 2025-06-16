package service

import (
	"encoding/base64"
	"fmt"
	"strconv"
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

	// Генерируем уникальный номер документа
	var uniqueNumber string
	var attempts int
	for attempts = 0; attempts < 10; attempts++ {
		proposedNumber := s.generateDocumentNumber(doc.Type)
		var count int64
		err := s.repo.GetDB().Model(&model.Document{}).Where("number = ?", proposedNumber).Count(&count).Error
		if err == nil && count == 0 {
			uniqueNumber = proposedNumber
			break
		}
	}
	if attempts == 10 {
		return &model.DocumentResponse{
			Message: "Не удалось сгенерировать уникальный номер документа",
		}
	}
	doc.Number = uniqueNumber

	// Устанавливаем статус черновика
	doc.Status = "draft"

	// Если дата не установлена, используем текущую
	if doc.Date.IsZero() {
		doc.Date = time.Now()
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

func (s *DocumentService) GetDocument(id uint) *model.DocumentResponse {
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

func (s *DocumentService) DeleteDocument(id uint) *model.DocumentResponse {
	response := s.repo.DeleteDocument(id)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) ApproveDocument(id uint, approvedByID uint) *model.DocumentResponse {
	response := s.repo.ApproveDocument(id, approvedByID)
	return &model.DocumentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *DocumentService) ExportDocument(id uint) *model.DocumentExportResponse {
	// старый экспорт в удобный для нас формат
	// Создаем сервис экспорта
	exportService := NewExportService(s)
	content, err := exportService.ExportDocument(id)
	if err != nil {
		return &model.DocumentExportResponse{
			Message: fmt.Sprintf("Ошибка экспорта документа: %v", err),
		}
	}

	// Преобразуем байты в строку для отправки в фронтенд
	return &model.DocumentExportResponse{
		Content: base64.StdEncoding.EncodeToString(content),
		Message: "Документ успешно экспортирован",
	}
}

func (s *DocumentService) ExportDocumentGOST(id uint) *model.DocumentExportResponse {
    exportService := NewExportService(s)
    content, err := exportService.ExportDocumentGOST(id)
    if err != nil {
        return &model.DocumentExportResponse{Message: fmt.Sprintf("Ошибка экспорта (ГОСТ): %v", err)}
    }
    return &model.DocumentExportResponse{Content: base64.StdEncoding.EncodeToString(content), Message: "Документ ГОСТ успешно экспортирован"}
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

	// Получаем максимальный номер для текущего года и типа документа
	var maxNumber string
	pattern := fmt.Sprintf("%s-%d-%%", prefix, year)
	err := s.repo.GetDB().Model(&model.Document{}).
		Where("type = ? AND number LIKE ?", docType, pattern).
		Select("number").
		Order("number DESC").
		Limit(1).
		Scan(&maxNumber).Error

	if err != nil {
		// Если нет документов для этого года, начинаем с 1
		return fmt.Sprintf("%s-%d-001", prefix, year)
	}

	// Извлекаем номер из последнего документа
	lastNumber := "001"
	if len(maxNumber) >= 11 {
		lastNumber = maxNumber[8:11]
	}

	// Преобразуем номер в число и увеличиваем на 1
	count, err := strconv.Atoi(lastNumber)
	if err != nil {
		count = 1
	}

	return fmt.Sprintf("%s-%d-%03d", prefix, year, count+1)
}
