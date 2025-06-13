package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) CreateLocation(location *model.Location) model.Response[*model.Location] {
	result := r.db.Create(location)
	if result.Error != nil {
		return model.Response[*model.Location]{
			Message: "Ошибка при создании местоположения",
		}
	}

	return model.Response[*model.Location]{
		Message: "Местоположение успешно создано",
		Model:   location,
	}
}

func (r *LocationRepository) GetLocation(id int) model.Response[*model.Location] {
	var location model.Location
	result := r.db.First(&location, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Response[*model.Location]{
				Message: "Местоположение не найдено",
			}
		}
		return model.Response[*model.Location]{
			Message: "Ошибка при получении местоположения",
		}
	}

	return model.Response[*model.Location]{
		Message: "Местоположение успешно получено",
		Model:   &location,
	}
}

func (r *LocationRepository) GetAllLocations() model.Response[[]model.Location] {
	var locations []model.Location
	result := r.db.Find(&locations)
	if result.Error != nil {
		return model.Response[[]model.Location]{
			Message: "Ошибка при получении списка местоположений",
		}
	}

	return model.Response[[]model.Location]{
		Message: "Список местоположений успешно получен",
		Model:   locations,
	}
}

func (r *LocationRepository) UpdateLocation(location *model.Location) model.Response[*model.Location] {
	result := r.db.Save(location)
	if result.Error != nil {
		return model.Response[*model.Location]{
			Message: "Ошибка при обновлении местоположения",
		}
	}

	return model.Response[*model.Location]{
		Message: "Местоположение успешно обновлено",
		Model:   location,
	}
}

func (r *LocationRepository) DeleteLocation(id int) model.Response[*model.Location] {
	var location model.Location
	result := r.db.First(&location, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Response[*model.Location]{
				Message: "Местоположение не найдено",
			}
		}
		return model.Response[*model.Location]{
			Message: "Ошибка при получении местоположения",
		}
	}

	result = r.db.Delete(&location)
	if result.Error != nil {
		return model.Response[*model.Location]{
			Message: "Ошибка при удалении местоположения",
		}
	}

	return model.Response[*model.Location]{
		Message: "Местоположение успешно удалено",
		Model:   &location,
	}
}

func (r *LocationRepository) GetLocationByEquipment(equipmentID int) model.Response[[]model.Location] {
	var locations []model.Location
	result := r.db.Joins("JOIN equipment ON equipment.location_id = locations.id").
		Where("equipment.id = ?", equipmentID).
		Find(&locations)
	if result.Error != nil {
		return model.Response[[]model.Location]{
			Message: "Ошибка при получении местоположений по оборудованию",
		}
	}

	return model.Response[[]model.Location]{
		Message: "Местоположения по оборудованию успешно получены",
		Model:   locations,
	}
}
