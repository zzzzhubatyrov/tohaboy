package service

import (
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
)

type AuthServiceInterface interface {
	Login(user map[string]string) (*model.LoginResponse, error)
	Register(user map[string]string) (*model.User, error)
}

type UserServiceInterface interface {
	GetUser(username string) model.Response[*model.User]
}

type EquipmentServiceInterface interface {
	CreateEquipment(equipment *model.Equipment) *model.EquipmentResponse
	GetEquipment(id int) *model.EquipmentResponse
	GetAllEquipment() *model.EquipmentListResponse
	UpdateEquipment(equipment *model.Equipment) *model.EquipmentResponse
	DeleteEquipment(id int) *model.EquipmentResponse
	GetEquipmentByLocation(locationID int) *model.EquipmentListResponse
	GetEquipmentBySupplier(supplierID int) *model.EquipmentListResponse
}

type SupplierServiceInterface interface {
	CreateSupplier(supplier *model.Supplier) *model.SupplierResponse
	GetSupplier(id int) *model.SupplierResponse
	GetAllSuppliers() *model.SupplierListResponse
	UpdateSupplier(supplier *model.Supplier) *model.SupplierResponse
	DeleteSupplier(id int) *model.SupplierResponse
	GetSupplierByEquipment(equipmentID int) *model.SupplierListResponse
}

type LocationServiceInterface interface {
	CreateLocation(location *model.Location) *model.LocationResponse
	GetLocation(id int) *model.LocationResponse
	GetAllLocations() *model.LocationListResponse
	UpdateLocation(location *model.Location) *model.LocationResponse
	DeleteLocation(id int) *model.LocationResponse
	GetLocationByEquipment(equipmentID int) *model.LocationListResponse
}

type MovementServiceInterface interface {
	CreateMovement(movement *model.Movement) *model.MovementResponse
	GetMovement(id int) *model.MovementResponse
	GetAllMovements() *model.MovementListResponse
	GetMovementsByEquipment(equipmentID int) *model.MovementListResponse
	GetMovementsByLocation(locationID int) *model.MovementListResponse
	UpdateMovement(movement *model.Movement) *model.MovementResponse
	DeleteMovement(id int) *model.MovementResponse
}

type DocumentServiceInterface interface {
	CreateDocument(doc *model.Document) *model.DocumentResponse
	GetDocument(id int) *model.DocumentResponse
	GetAllDocuments() *model.DocumentListResponse
	UpdateDocument(doc *model.Document) *model.DocumentResponse
	DeleteDocument(id int) *model.DocumentResponse
	ApproveDocument(id int, approvedByID int) *model.DocumentResponse
}

type Service struct {
	AuthServiceInterface
	UserService      UserServiceInterface
	EquipmentService EquipmentServiceInterface
	SupplierService  SupplierServiceInterface
	LocationService  LocationServiceInterface
	MovementService  MovementServiceInterface
	DocumentService  DocumentServiceInterface
}

func NewService(repos *repository.Repository) *Service {
	docService := NewDocumentService(repos.Document)
	return &Service{
		AuthServiceInterface: NewAuthService(repos.AuthRepositoryInterface),
		UserService:          NewUserService(repos.User),
		EquipmentService:     NewEquipmentService(repos.Equipment),
		SupplierService:      NewSupplierService(repos.Supplier),
		LocationService:      NewLocationService(repos.Location),
		MovementService:      NewMovementService(repos.Movement, docService),
		DocumentService:      docService,
	}
}
