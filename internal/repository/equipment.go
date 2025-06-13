package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type EquipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) CreateEquipment(equipment *model.Equipment) model.Response[*model.Equipment] {
	if err := r.db.Create(equipment).Error; err != nil {
		return model.Response[*model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[*model.Equipment]{
		Model:   equipment,
		Message: "Оборудование успешно добавлено",
	}
}

func (r *EquipmentRepository) GetEquipment(id int) model.Response[*model.Equipment] {
	var equipment model.Equipment

	if err := r.db.Preload("Location").
		Preload("Supplier").
		Preload("Movements").
		First(&equipment, id).Error; err != nil {
		return model.Response[*model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[*model.Equipment]{
		Model:   &equipment,
		Message: "Оборудование найдено",
	}
}

func (r *EquipmentRepository) UpdateEquipment(equipment *model.Equipment) model.Response[*model.Equipment] {
	if err := r.db.Save(equipment).Error; err != nil {
		return model.Response[*model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[*model.Equipment]{
		Model:   equipment,
		Message: "Оборудование обновлено",
	}
}

func (r *EquipmentRepository) DeleteEquipment(id int) model.Response[*model.Equipment] {
	var equipment model.Equipment
	if err := r.db.First(&equipment, id).Error; err != nil {
		return model.Response[*model.Equipment]{
			Model:   nil,
			Message: "Оборудование не найдено: " + err.Error(),
		}
	}

	if err := r.db.Delete(&equipment).Error; err != nil {
		return model.Response[*model.Equipment]{
			Model:   nil,
			Message: "Ошибка при удалении: " + err.Error(),
		}
	}

	return model.Response[*model.Equipment]{
		Model:   &equipment,
		Message: "Оборудование удалено",
	}
}

func (r *EquipmentRepository) GetAllEquipment() model.Response[[]model.Equipment] {
	var equipment []model.Equipment

	if err := r.db.Preload("Location").
		Preload("Supplier").
		Preload("Movements").
		Find(&equipment).Error; err != nil {
		return model.Response[[]model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Equipment]{
		Model:   equipment,
		Message: "Все оборудование загружено",
	}
}

func (r *EquipmentRepository) GetEquipmentByLocation(locationID int) model.Response[[]model.Equipment] {
	var equipment []model.Equipment

	if err := r.db.Where("location_id = ?", locationID).
		Preload("Location").Preload("Supplier").Preload("Movements").
		Find(&equipment).Error; err != nil {
		return model.Response[[]model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Equipment]{
		Model:   equipment,
		Message: "Оборудование по локации загружено",
	}
}

func (r *EquipmentRepository) GetEquipmentBySupplier(supplierID int) model.Response[[]model.Equipment] {
	var equipment []model.Equipment

	if err := r.db.Where("supplier_id = ?", supplierID).
		Preload("Location").Preload("Supplier").Preload("Movements").
		Find(&equipment).Error; err != nil {
		return model.Response[[]model.Equipment]{
			Model:   nil,
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Equipment]{
		Model:   equipment,
		Message: "Оборудование по поставщику загружено",
	}
}
