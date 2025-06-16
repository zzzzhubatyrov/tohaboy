package service

import (
	"fmt"
	"time"
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type MovementService struct {
	repo repository.MovementRepositoryInterface
	docs DocumentServiceInterface
}

func NewMovementService(repo repository.MovementRepositoryInterface, docs DocumentServiceInterface) *MovementService {
	return &MovementService{repo: repo, docs: docs}
}

func (s *MovementService) CreateMovement(movement *model.Movement) *model.MovementResponse {
	// Валидация
	if err := s.validateMovement(movement); err != nil {
		return &model.MovementResponse{
			Message: err.Error(),
		}
	}

	// Устанавливаем дату создания
	movement.Date = time.Now()

	// Создаем документ перемещения
	doc := &model.Document{
		Type:        "transfer",
		Date:        movement.Date,
		Status:      "completed",
		LocationID:  movement.ToLocationID,
		CreatedByID: movement.CreatedByID,
		Items: []model.DocumentItem{
			{
				EquipmentID: movement.EquipmentID,
				Quantity:    movement.Quantity,
			},
		},
	}

	// Создаем документ
	docResponse := s.docs.CreateDocument(doc)
	if docResponse.Model == nil {
		return &model.MovementResponse{
			Message: fmt.Sprintf("Ошибка создания документа: %s", docResponse.Message),
		}
	}

	// Связываем движение с созданным документом
	movement.DocumentID = docResponse.Model.ID

	// Создаем запись о перемещении
	response := s.repo.CreateMovement(movement)
	return &model.MovementResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) GetMovement(id uint) *model.MovementResponse {
	response := s.repo.GetMovement(id)
	return &model.MovementResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) GetAllMovements() *model.MovementListResponse {
	response := s.repo.GetAllMovements()
	return &model.MovementListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) UpdateMovement(movement *model.Movement) *model.MovementResponse {
	response := s.repo.UpdateMovement(movement)
	return &model.MovementResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) DeleteMovement(id uint) *model.MovementResponse {
	response := s.repo.DeleteMovement(id)
	return &model.MovementResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) GetMovementsByEquipment(equipmentID uint) *model.MovementListResponse {
	response := s.repo.GetMovementsByEquipment(equipmentID)
	return &model.MovementListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *MovementService) GetMovementsByLocation(locationID uint) *model.MovementListResponse {
	response := s.repo.GetMovementsByLocation(locationID)
	return &model.MovementListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

// Вспомогательные методы

func (s *MovementService) validateMovement(movement *model.Movement) error {
	if movement.EquipmentID == 0 {
		return fmt.Errorf("оборудование не указано")
	}

	if movement.FromLocationID == 0 {
		return fmt.Errorf("начальное местоположение не указано")
	}

	if movement.ToLocationID == 0 {
		return fmt.Errorf("конечное местоположение не указано")
	}

	if movement.FromLocationID == movement.ToLocationID {
		return fmt.Errorf("начальное и конечное местоположение совпадают")
	}

	if movement.Quantity <= 0 {
		return fmt.Errorf("количество должно быть больше нуля")
	}

	if movement.CreatedByID == 0 {
		return fmt.Errorf("пользователь не указан")
	}

	return nil
}
