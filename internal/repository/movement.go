package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type MovementRepository struct {
	db *gorm.DB
}

func NewMovementRepository(db *gorm.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

func (r *MovementRepository) CreateMovement(movement *model.Movement) model.Response[*model.Movement] {
	// Начинаем транзакцию
	tx := r.db.Begin()

	// Получаем оборудование для проверки и обновления
	var equipment model.Equipment
	if err := tx.First(&equipment, movement.EquipmentID).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Movement]{
			Message: "Оборудование не найдено",
		}
	}

	// Проверяем достаточность количества
	if movement.Quantity > equipment.Quantity {
		tx.Rollback()
		return model.Response[*model.Movement]{
			Message: "Недостаточное количество оборудования",
		}
	}

	// Создаем запись о перемещении
	if err := tx.Create(movement).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Movement]{
			Message: err.Error(),
		}
	}

	// Обновляем местоположение оборудования
	equipment.LocationID = movement.ToLocationID
	if err := tx.Save(&equipment).Error; err != nil {
		tx.Rollback()
		return model.Response[*model.Movement]{
			Message: err.Error(),
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return model.Response[*model.Movement]{
			Message: err.Error(),
		}
	}

	// Загружаем созданное перемещение со всеми связями
	if err := r.db.Preload("Equipment").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("CreatedBy").
		First(movement, movement.ID).Error; err != nil {
		return model.Response[*model.Movement]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Movement]{
		Model: movement,
	}
}

func (r *MovementRepository) GetMovement(id uint) model.Response[*model.Movement] {
	var movement model.Movement

	if err := r.db.Preload("Equipment").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("CreatedBy").
		First(&movement, id).Error; err != nil {
		return model.Response[*model.Movement]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Movement]{
		Model: &movement,
	}
}

func (r *MovementRepository) UpdateMovement(equipmentMovement *model.Movement) model.Response[*model.Movement] {
	if err := r.db.Save(equipmentMovement).Error; err != nil {
		return model.Response[*model.Movement]{}
	}

	return model.Response[*model.Movement]{}
}

func (r *MovementRepository) DeleteMovement(id uint) model.Response[*model.Movement] {
	if err := r.db.Delete(&model.Movement{}, id).Error; err != nil {
		return model.Response[*model.Movement]{}
	}

	return model.Response[*model.Movement]{}
}

func (r *MovementRepository) GetAllMovements() model.Response[[]model.Movement] {
	var movements []model.Movement

	if err := r.db.Preload("Equipment").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("CreatedBy").
		Order("date DESC").
		Find(&movements).Error; err != nil {
		return model.Response[[]model.Movement]{
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Movement]{
		Model: movements,
	}
}

func (r *MovementRepository) GetMovementsByEquipment(equipmentID uint) model.Response[[]model.Movement] {
	var movements []model.Movement

	if err := r.db.Preload("Equipment").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("CreatedBy").
		Where("equipment_id = ?", equipmentID).
		Order("date DESC").
		Find(&movements).Error; err != nil {
		return model.Response[[]model.Movement]{
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Movement]{
		Model: movements,
	}
}

func (r *MovementRepository) GetMovementsByLocation(locationID uint) model.Response[[]model.Movement] {
	var movements []model.Movement

	if err := r.db.Preload("Equipment").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("CreatedBy").
		Where("from_location_id = ? OR to_location_id = ?", locationID, locationID).
		Order("date DESC").
		Find(&movements).Error; err != nil {
		return model.Response[[]model.Movement]{
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Movement]{
		Model: movements,
	}
}
