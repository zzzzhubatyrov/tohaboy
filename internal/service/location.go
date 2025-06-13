package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type LocationService struct {
	repo repository.LocationRepositoryInterface
}

func NewLocationService(repo repository.LocationRepositoryInterface) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) CreateLocation(location *model.Location) *model.LocationResponse {
	response := s.repo.CreateLocation(location)
	return &model.LocationResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *LocationService) GetLocation(id int) *model.LocationResponse {
	response := s.repo.GetLocation(id)
	return &model.LocationResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *LocationService) GetAllLocations() *model.LocationListResponse {
	response := s.repo.GetAllLocations()
	return &model.LocationListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *LocationService) UpdateLocation(location *model.Location) *model.LocationResponse {
	response := s.repo.UpdateLocation(location)
	return &model.LocationResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *LocationService) DeleteLocation(id int) *model.LocationResponse {
	response := s.repo.DeleteLocation(id)
	return &model.LocationResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}

func (s *LocationService) GetLocationByEquipment(equipmentID int) *model.LocationListResponse {
	response := s.repo.GetLocationByEquipment(equipmentID)
	return &model.LocationListResponse{
		Model:   response.Model,
		Message: response.Message,
	}
}
