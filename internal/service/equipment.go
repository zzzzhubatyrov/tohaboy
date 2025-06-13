package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type EquipmentService struct {
	repo repository.EquipmentRepositoryInterface
}

func NewEquipmentService(repo repository.EquipmentRepositoryInterface) *EquipmentService {
	return &EquipmentService{repo: repo}
}

func (s *EquipmentService) CreateEquipment(equipment *model.Equipment) *model.EquipmentResponse {
	response := s.repo.CreateEquipment(equipment)
	return &model.EquipmentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) GetEquipment(id int) *model.EquipmentResponse {
	response := s.repo.GetEquipment(id)
	return &model.EquipmentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) GetAllEquipment() *model.EquipmentListResponse {
	response := s.repo.GetAllEquipment()
	return &model.EquipmentListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) UpdateEquipment(equipment *model.Equipment) *model.EquipmentResponse {
	response := s.repo.UpdateEquipment(equipment)
	return &model.EquipmentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) DeleteEquipment(id int) *model.EquipmentResponse {
	response := s.repo.DeleteEquipment(id)
	return &model.EquipmentResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) GetEquipmentByLocation(locationID int) *model.EquipmentListResponse {
	response := s.repo.GetEquipmentByLocation(locationID)
	return &model.EquipmentListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *EquipmentService) GetEquipmentBySupplier(supplierID int) *model.EquipmentListResponse {
	response := s.repo.GetEquipmentBySupplier(supplierID)
	return &model.EquipmentListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}
