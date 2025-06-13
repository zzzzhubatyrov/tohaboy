package repository

import (
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Login(user *model.User) (*model.User, error)
	Register(user *model.User) (*model.User, error)
}

type UserRepositoryInterface interface {
	GetUser(username string) model.Response[*model.User]
	GetByID(id int) model.Response[*model.User]
	Update(user *model.User) model.Response[*model.User]
	Delete(id int) model.Response[*model.User]
}

type EquipmentRepositoryInterface interface {
	CreateEquipment(equipment *model.Equipment) model.Response[*model.Equipment]
	GetEquipment(id int) model.Response[*model.Equipment]
	GetAllEquipment() model.Response[[]model.Equipment]
	UpdateEquipment(equipment *model.Equipment) model.Response[*model.Equipment]
	DeleteEquipment(id int) model.Response[*model.Equipment]
	GetEquipmentByLocation(locationID int) model.Response[[]model.Equipment]
	GetEquipmentBySupplier(supplierID int) model.Response[[]model.Equipment]
}

type SupplierRepositoryInterface interface {
	CreateSupplier(supplier *model.Supplier) model.Response[*model.Supplier]
	GetSupplier(id int) model.Response[*model.Supplier]
	GetAllSuppliers() model.Response[[]model.Supplier]
	UpdateSupplier(supplier *model.Supplier) model.Response[*model.Supplier]
	DeleteSupplier(id int) model.Response[*model.Supplier]
	GetSupplierByEquipment(equipmentID int) model.Response[[]model.Supplier]
}

type LocationRepositoryInterface interface {
	CreateLocation(location *model.Location) model.Response[*model.Location]
	GetLocation(id int) model.Response[*model.Location]
	GetAllLocations() model.Response[[]model.Location]
	UpdateLocation(location *model.Location) model.Response[*model.Location]
	DeleteLocation(id int) model.Response[*model.Location]
	GetLocationByEquipment(equipmentID int) model.Response[[]model.Location]
}

type MovementRepositoryInterface interface {
	CreateMovement(movement *model.Movement) model.Response[*model.Movement]
	GetMovement(id int) model.Response[*model.Movement]
	GetAllMovements() model.Response[[]model.Movement]
	GetMovementsByEquipment(equipmentID int) model.Response[[]model.Movement]
	GetMovementsByLocation(locationID int) model.Response[[]model.Movement]
	UpdateMovement(movement *model.Movement) model.Response[*model.Movement]
	DeleteMovement(id int) model.Response[*model.Movement]
}

type DocumentRepositoryInterface interface {
	CreateDocument(doc *model.Document) model.Response[*model.Document]
	GetDocument(id int) model.Response[*model.Document]
	GetAllDocuments() model.Response[[]model.Document]
	UpdateDocument(doc *model.Document) model.Response[*model.Document]
	DeleteDocument(id int) model.Response[*model.Document]
	ApproveDocument(id int, approvedByID int) model.Response[*model.Document]
}

type Repository struct {
	AuthRepositoryInterface
	User      UserRepositoryInterface
	Equipment EquipmentRepositoryInterface
	Supplier  SupplierRepositoryInterface
	Location  LocationRepositoryInterface
	Movement  MovementRepositoryInterface
	Document  DocumentRepositoryInterface
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AuthRepositoryInterface: NewAuthRepo(db),
		User:                    NewUserRepository(db),
		Equipment:               NewEquipmentRepository(db),
		Supplier:                NewSupplierRepository(db),
		Location:                NewLocationRepository(db),
		Movement:                NewMovementRepository(db),
		Document:                NewDocumentRepository(db),
	}
}
