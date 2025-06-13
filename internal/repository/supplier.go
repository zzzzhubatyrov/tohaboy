package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) CreateSupplier(supplier *model.Supplier) model.Response[*model.Supplier] {
	if err := r.db.Create(supplier).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Supplier]{
		Model: supplier,
	}
}

func (r *SupplierRepository) GetSupplier(id int) model.Response[*model.Supplier] {
	var supplier model.Supplier

	if err := r.db.Preload("Equipment").First(&supplier, id).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Supplier]{
		Model: &supplier,
	}
}

func (r *SupplierRepository) UpdateSupplier(supplier *model.Supplier) model.Response[*model.Supplier] {
	// Сначала проверяем существование
	var existingSupplier model.Supplier
	if err := r.db.First(&existingSupplier, supplier.ID).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: "Поставщик не найден",
		}
	}

	if err := r.db.Save(supplier).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	// Получаем обновленные данные
	if err := r.db.Preload("Equipment").First(&existingSupplier, supplier.ID).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Supplier]{
		Model: &existingSupplier,
	}
}

func (r *SupplierRepository) DeleteSupplier(id int) model.Response[*model.Supplier] {
	// Сначала получаем поставщика для возврата данных
	var supplier model.Supplier
	if err := r.db.Preload("Equipment").First(&supplier, id).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	// Затем удаляем
	if err := r.db.Delete(&supplier).Error; err != nil {
		return model.Response[*model.Supplier]{
			Message: err.Error(),
		}
	}

	return model.Response[*model.Supplier]{
		Model: &supplier,
	}
}

func (r *SupplierRepository) GetAllSuppliers() model.Response[[]model.Supplier] {
	var suppliers []model.Supplier

	if err := r.db.Preload("Equipment").Find(&suppliers).Error; err != nil {
		return model.Response[[]model.Supplier]{
			Message: err.Error(),
		}
	}

	return model.Response[[]model.Supplier]{
		Model: suppliers,
	}
}

func (r *SupplierRepository) GetSupplierByEquipment(equipmentID int) model.Response[[]model.Supplier] {
	var suppliers []model.Supplier
	result := r.db.Joins("JOIN equipment ON equipment.supplier_id = suppliers.id").
		Where("equipment.id = ?", equipmentID).
		Find(&suppliers)
	if result.Error != nil {
		return model.Response[[]model.Supplier]{
			Message: "Ошибка при получении поставщиков по оборудованию",
		}
	}

	return model.Response[[]model.Supplier]{
		Message: "Поставщики по оборудованию успешно получены",
		Model:   suppliers,
	}
}
