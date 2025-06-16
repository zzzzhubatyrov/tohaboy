package main

import (
	"embed"
	"fmt"
	"math/rand"
	"time"
	"tohaboy/internal/data"
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
	"tohaboy/internal/service"
	"tohaboy/internal/storage"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var err error

	// Create an instance of the app structure
	app := NewApp()

	// Create DB connection
	db := storage.NewStorage("invent.db")

	models := []interface{}{
		&model.User{},
		&model.Equipment{},
		&model.Supplier{},
		&model.Location{},
		&model.Movement{},
		&model.Document{},
		&model.DocumentItem{},
		&model.Category{},
	}

	// Only drop tables if explicitly needed for development
	if true { // Change to true only when needed during development
		if err = db.DropTables(models); err != nil {
			panic(err)
		}
	}

	// Migrate database
	if err = db.Migrate(models); err != nil {
		panic(err)
	}

	// Generate test data only if tables are empty
	if err = db.GetDB().First(&model.User{}).Error; err != nil {
		if err = generateTestData(db.GetDB()); err != nil {
			panic(fmt.Sprintf("Error generating test data: %v", err))
		}
	}

	// Create services
	svc := service.NewService(repository.NewRepository(db.GetDB()))

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Инвентаризация и управление оборудованием",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			svc.AuthServiceInterface,
			svc.UserService,
			svc.EquipmentService,
			svc.SupplierService,
			svc.LocationService,
			svc.MovementService,
			svc.DocumentService,
			svc.CategoryService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func generateTestData(db *gorm.DB) error {
	// Create admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
	}
	if err := db.Create(admin).Error; err != nil {
		return err
	}

	// Заполняем базу тестовыми данными
	if err := data.SeedDatabase(db); err != nil {
		return fmt.Errorf("error seeding database: %v", err)
	}

	return nil
}

func generateSerialNumber() string {
	return fmt.Sprintf("SN-%s-%d",
		time.Now().Format("20060102"),
		rand.Intn(1000),
	)
}

func generateDocNumber(docType string, count int) string {
	return fmt.Sprintf("%s-%s-%03d",
		docType,
		time.Now().Format("20060102"),
		count,
	)
}
