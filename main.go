package main

import (
	"embed"
	"fmt"
	"math/rand"
	"time"
	"tohaboy/internal/model"
	"tohaboy/internal/repository"
	"tohaboy/internal/service"
	"tohaboy/internal/storage"

	"github.com/brianvoe/gofakeit/v7"
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
	}

	// Удаляем старые таблицы
	if err = db.DropTables(models); err != nil {
		panic(err)
	}

	// Создаем новые таблицы
	if err = db.Migrate(models); err != nil {
		panic(err)
	}

	// Генерируем тестовые данные
	if err = generateTestData(db.GetDB()); err != nil {
		panic(fmt.Sprintf("Ошибка при генерации тестовых данных: %v", err))
	}

	// Create services
	svc := service.NewService(repository.NewRepository(db.GetDB()))

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "tohaboy",
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
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func SeedTestData(db *gorm.DB) error {
	// Создаем 10 поставщиков
	suppliers := make([]model.Supplier, 10)
	for i := 0; i < 10; i++ {
		suppliers[i] = model.Supplier{
			Name: fmt.Sprintf("Поставщик %d", i+1),
		}
	}
	if err := db.Create(&suppliers).Error; err != nil {
		return fmt.Errorf("ошибка при создании поставщиков: %v", err)
	}

	// Создаем 10 локаций
	locations := make([]model.Location, 10)
	for i := 0; i < 10; i++ {
		locations[i] = model.Location{
			Name:        fmt.Sprintf("Локация %d", i+1),
			Description: fmt.Sprintf("Описание локации %d", i+1),
		}
	}
	if err := db.Create(&locations).Error; err != nil {
		return fmt.Errorf("ошибка при создании локаций: %v", err)
	}

	// Создаем 10 единиц оборудования
	equipments := make([]model.Equipment, 10)
	for i := 0; i < 10; i++ {
		equipments[i] = model.Equipment{
			Name:         fmt.Sprintf("Оборудование %d", i+1),
			SerialNumber: gofakeit.UUID(),
			Category:     gofakeit.RandomString([]string{"PC", "Router", "Chair"}),
			Description:  fmt.Sprintf("Описание оборудования %d", i+1),
			Price:        float64(rand.Intn(10000) + 1000), // Случайная цена от 1000 до 11000
			Status:       []string{"active", "inactive", "maintenance"}[rand.Intn(3)],
			LocationID:   locations[rand.Intn(len(locations))].ID,
			SupplierID:   suppliers[rand.Intn(len(suppliers))].ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}
	if err := db.Create(&equipments).Error; err != nil {
		return fmt.Errorf("ошибка при создании оборудования: %v", err)
	}

	// Создаем тестового пользователя для перемещений
	location := model.Location{ID: 11} // Предполагаем, что у вас есть пользователь с ID=1
	if err := db.FirstOrCreate(&location).Error; err != nil {
		return fmt.Errorf("ошибка при создании тестового пользователя: %v", err)
	}

	// Создаем 10 записей о перемещении оборудования
	movements := make([]model.Movement, 10)
	for i := 0; i < 10; i++ {
		movements[i] = model.Movement{
			EquipmentID:    equipments[rand.Intn(len(equipments))].ID,
			FromLocationID: locations[rand.Intn(len(locations))].ID,
			ToLocationID:   locations[rand.Intn(len(locations))].ID,
			// Date:           model.Date(time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour)), // Дата от сегодня до 30 дней назад
		}
	}
	if err := db.Create(&movements).Error; err != nil {
		return fmt.Errorf("ошибка при создании записей о перемещении: %v", err)
	}

	return nil
}

func generateTestData(db *gorm.DB) error {
	// Генерация пользователей
	users := []model.User{
		{Username: "admin", Role: "admin"},
		{Username: "manager", Role: "manager"},
		{Username: "auditor", Role: "auditor"},
	}

	// Хешируем пароль "admin" для всех пользователей
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	for i := range users {
		// Проверяем существование пользователя
		var existingUser model.User
		if err := db.Where("username = ?", users[i].Username).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Пользователь не существует, создаем нового
				users[i].Password = string(hashedPassword)
				if err := db.Create(&users[i]).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Проверяем наличие поставщиков
	var supplierCount int64
	db.Model(&model.Supplier{}).Count(&supplierCount)
	if supplierCount == 0 {
		// Генерация поставщиков
		suppliers := []model.Supplier{
			{Name: "ООО Техноснаб", ContactInfo: "Иванов Иван, +7 (999) 123-45-67, ivanov@technosnab.ru"},
			{Name: "Компьютерный Мир", ContactInfo: "Петров Петр, +7 (999) 234-56-78, petrov@compworld.ru"},
			{Name: "Офисные Решения", ContactInfo: "Сидорова Анна, +7 (999) 345-67-89, sidorova@officesolutions.ru"},
			{Name: "МебельПром", ContactInfo: "Козлов Дмитрий, +7 (999) 456-78-90, kozlov@mebelprom.ru"},
			{Name: "ЭлектроТрейд", ContactInfo: "Морозова Елена, +7 (999) 567-89-01, morozova@electrotrade.ru"},
		}

		for _, supplier := range suppliers {
			if err := db.Create(&supplier).Error; err != nil {
				return err
			}
		}
	}

	// Проверяем наличие местоположений
	var locationCount int64
	db.Model(&model.Location{}).Count(&locationCount)
	if locationCount == 0 {
		// Генерация местоположений
		locations := []model.Location{
			{Name: "Главный офис", Description: "Основное здание, 1 этаж", Address: "ул. Ленина, 1"},
			{Name: "Склад №1", Description: "Складское помещение", Address: "ул. Складская, 10"},
			{Name: "Отдел разработки", Description: "Офис разработчиков, 3 этаж", Address: "ул. Ленина, 1"},
			{Name: "Серверная", Description: "Помещение с серверами", Address: "ул. Ленина, 1"},
			{Name: "Бухгалтерия", Description: "Бухгалтерия, 2 этаж", Address: "ул. Ленина, 1"},
		}

		for _, location := range locations {
			if err := db.Create(&location).Error; err != nil {
				return err
			}
		}
	}

	// Проверяем наличие оборудования
	var equipmentCount int64
	db.Model(&model.Equipment{}).Count(&equipmentCount)
	if equipmentCount == 0 {
		// Получаем существующие местоположения и поставщиков
		var locations []model.Location
		var suppliers []model.Supplier
		if err := db.Find(&locations).Error; err != nil {
			return err
		}
		if err := db.Find(&suppliers).Error; err != nil {
			return err
		}

		// Генерация оборудования
		equipment := []model.Equipment{
			{
				Name:         "Ноутбук Dell XPS 15",
				SerialNumber: generateSerialNumber(),
				Category:     "Компьютеры",
				Description:  "Мощный ноутбук для разработки",
				Price:        120000,
				Status:       "available",
				Quantity:     5,
				LocationID:   locations[0].ID,
				SupplierID:   suppliers[0].ID,
			},
			{
				Name:         "Монитор LG 27\"",
				SerialNumber: generateSerialNumber(),
				Category:     "Мониторы",
				Description:  "27-дюймовый монитор 4K",
				Price:        35000,
				Status:       "available",
				Quantity:     10,
				LocationID:   locations[1].ID,
				SupplierID:   suppliers[1].ID,
			},
			{
				Name:         "Кресло офисное",
				SerialNumber: generateSerialNumber(),
				Category:     "Мебель",
				Description:  "Эргономичное офисное кресло",
				Price:        15000,
				Status:       "available",
				Quantity:     20,
				LocationID:   locations[2].ID,
				SupplierID:   suppliers[3].ID,
			},
			{
				Name:         "Сервер Dell PowerEdge",
				SerialNumber: generateSerialNumber(),
				Category:     "Серверы",
				Description:  "Сервер для разработки",
				Price:        450000,
				Status:       "in_use",
				Quantity:     2,
				LocationID:   locations[3].ID,
				SupplierID:   suppliers[0].ID,
			},
			{
				Name:         "МФУ HP LaserJet",
				SerialNumber: generateSerialNumber(),
				Category:     "Принтеры",
				Description:  "Многофункциональное устройство",
				Price:        45000,
				Status:       "available",
				Quantity:     3,
				LocationID:   locations[4].ID,
				SupplierID:   suppliers[1].ID,
			},
		}

		for _, eq := range equipment {
			if err := db.Create(&eq).Error; err != nil {
				return err
			}
		}
	}

	// Проверяем наличие документов
	var documentCount int64
	db.Model(&model.Document{}).Count(&documentCount)
	if documentCount == 0 {
		// Получаем пользователя admin
		var adminUser model.User
		if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
			return err
		}

		// Получаем оборудование
		var equipment []model.Equipment
		if err := db.Find(&equipment).Error; err != nil {
			return err
		}

		// Генерация документов
		documents := []model.Document{
			{
				Type:        "inventory",
				Number:      generateDocNumber("INV", 1),
				Status:      "completed",
				Date:        model.Date(time.Now().Add(-24 * time.Hour)),
				CreatedByID: adminUser.ID,
				LocationID:  equipment[0].LocationID,
				Items: []model.DocumentItem{
					{
						EquipmentID:    equipment[0].ID,
						Quantity:       2,
						ActualQuantity: 2,
						Price:          equipment[0].Price,
						TotalPrice:     equipment[0].Price * 2,
					},
				},
			},
			{
				Type:        "transfer",
				Number:      generateDocNumber("TRN", 1),
				Status:      "completed",
				Date:        model.Date(time.Now().Add(-48 * time.Hour)),
				CreatedByID: adminUser.ID,
				LocationID:  equipment[1].LocationID,
				Items: []model.DocumentItem{
					{
						EquipmentID: equipment[1].ID,
						Quantity:    1,
						Price:       equipment[1].Price,
						TotalPrice:  equipment[1].Price,
					},
				},
			},
		}

		for _, doc := range documents {
			if err := db.Create(&doc).Error; err != nil {
				return err
			}
		}
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
