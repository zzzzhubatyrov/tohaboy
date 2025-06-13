package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type SupplierService struct {
	repo repository.SupplierRepositoryInterface
}

func NewSupplierService(repo repository.SupplierRepositoryInterface) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) CreateSupplier(supplier *model.Supplier) *model.SupplierResponse {
	response := s.repo.CreateSupplier(supplier)
	return &model.SupplierResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *SupplierService) GetSupplier(id int) *model.SupplierResponse {
	response := s.repo.GetSupplier(id)
	return &model.SupplierResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *SupplierService) GetAllSuppliers() *model.SupplierListResponse {
	response := s.repo.GetAllSuppliers()
	return &model.SupplierListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *SupplierService) UpdateSupplier(supplier *model.Supplier) *model.SupplierResponse {
	response := s.repo.UpdateSupplier(supplier)
	return &model.SupplierResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *SupplierService) DeleteSupplier(id int) *model.SupplierResponse {
	response := s.repo.DeleteSupplier(id)
	return &model.SupplierResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *SupplierService) GetSupplierByEquipment(equipmentID int) *model.SupplierListResponse {
	response := s.repo.GetSupplierByEquipment(equipmentID)
	return &model.SupplierListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}
