package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"tohaboy/internal/model"
	"tohaboy/internal/storage"
)

var (
	companyNames = []string{
		"ТехноСервис", "ИнновацияПлюс", "ПромОборудование", "МедТехника", "СтройМаш",
		"ЭлектроСистемы", "АвтоТехСнаб", "ПромСнабжение", "МедикалГрупп", "ТехноМир",
		"СтанкоТрейд", "ЭнергоТех", "МашСервис", "ПромКомплект", "МедОборудование",
	}

	streets = []string{
		"ул. Ленина", "пр. Мира", "ул. Гагарина", "ул. Пушкина", "пр. Победы",
		"ул. Советская", "ул. Московская", "пр. Ленинградский", "ул. Строителей",
	}

	cities = []string{
		"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань",
		"Нижний Новгород", "Челябинск", "Самара", "Уфа", "Ростов-на-Дону",
	}

	equipmentTypes = []string{
		"Станок токарный", "Фрезерный станок", "Сварочный аппарат", "Компрессор",
		"Генератор", "Погрузчик", "Конвейер", "Пресс гидравлический", "Дрель промышленная",
		"Шлифовальная машина", "Упаковочная машина", "Термопластавтомат",
	}

	equipmentBrands = []string{
		"HAAS", "DMG MORI", "Siemens", "ABB", "Fanuc", "Mitsubishi", "Bosch",
		"Schneider Electric", "Yamazaki Mazak", "TRUMPF", "Amada", "Okuma",
	}
)

func generatePhone() string {
	return fmt.Sprintf("+7 (%d) %d-%d-%d",
		rand.Intn(900)+100,
		rand.Intn(900)+100,
		rand.Intn(90)+10,
		rand.Intn(90)+10,
	)
}

func generateAddress() string {
	building := rand.Intn(200) + 1
	return fmt.Sprintf("%s, д. %d, %s", streets[rand.Intn(len(streets))], building, cities[rand.Intn(len(cities))])
}

func generateSerialNumber() string {
	return fmt.Sprintf("%s-%d-%d",
		string(rune('A'+rand.Intn(26))),
		rand.Intn(10000),
		time.Now().Year(),
	)
}

func generateSampleData() {
	rand.Seed(time.Now().UnixNano())

	// Initialize storage and migrate tables
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

	if err := db.Migrate(models); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Clear existing data
	if err := db.GetDB().Exec("DELETE FROM equipment").Error; err != nil {
		log.Printf("Error clearing equipment: %v", err)
	}
	if err := db.GetDB().Exec("DELETE FROM suppliers").Error; err != nil {
		log.Printf("Error clearing suppliers: %v", err)
	}

	// Генерируем 30 поставщиков
	for i := 0; i < 30; i++ {
		supplier := &model.Supplier{
			Name:        fmt.Sprintf("%s-%d", companyNames[rand.Intn(len(companyNames))], i+1),
			Description: fmt.Sprintf("Поставщик промышленного оборудования и комплектующих. Работаем с %d года.", time.Now().Year()-rand.Intn(20)),
			Address:     generateAddress(),
			Phone:       generatePhone(),
		}

		if err := db.GetDB().Create(supplier).Error; err != nil {
			log.Printf("Error creating supplier: %v", err)
			continue
		}

		// Генерируем от 2 до 5 единиц оборудования для каждого поставщика
		numEquipment := rand.Intn(4) + 2
		for j := 0; j < numEquipment; j++ {
			equipType := equipmentTypes[rand.Intn(len(equipmentTypes))]
			brand := equipmentBrands[rand.Intn(len(equipmentBrands))]
			equipment := &model.Equipment{
				Name:         fmt.Sprintf("%s %s", brand, equipType),
				SerialNumber: generateSerialNumber(),
				Description:  fmt.Sprintf("Промышленное оборудование %s производства %s", equipType, brand),
				SupplierID:   supplier.ID,
				Status:       "В эксплуатации",
				Quantity:     rand.Intn(5) + 1,
				Price:        float64(rand.Intn(1000000) + 100000),
			}

			if err := db.GetDB().Create(equipment).Error; err != nil {
				log.Printf("Error creating equipment: %v", err)
				continue
			}
		}
	}

	fmt.Println("Sample data generated successfully!")
}
