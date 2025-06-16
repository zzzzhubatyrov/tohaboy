package data

import (
	"time"
	"tohaboy/internal/model"

	"gorm.io/gorm"
)

// GetCategories возвращает список категорий оборудования
func GetCategories() []model.Category {
	return []model.Category{
		{ID: 1, Name: "Компьютерная техника", Description: "Компьютеры, ноутбуки, серверы и комплектующие"},
		{ID: 2, Name: "Сетевое оборудование", Description: "Маршрутизаторы, коммутаторы, точки доступа"},
		{ID: 3, Name: "Офисная техника", Description: "Принтеры, сканеры, МФУ"},
		{ID: 4, Name: "Мебель", Description: "Столы, стулья, шкафы"},
		{ID: 5, Name: "Инструменты", Description: "Ручной и электрический инструмент"},
		{ID: 6, Name: "Программное обеспечение", Description: "Лицензии на ПО"},
		{ID: 7, Name: "Телекоммуникационное оборудование", Description: "Телефоны, видеоконференцсвязь"},
		{ID: 8, Name: "Системы безопасности", Description: "Камеры, СКУД, сигнализация"},
	}
}

// GetLocations возвращает список местоположений
func GetLocations() []model.Location {
	return []model.Location{
		{ID: 1, Name: "Главный офис", Description: "Основное здание компании", Address: "ул. Ленина, 1"},
		{ID: 2, Name: "Серверная", Description: "Серверное помещение", Address: "ул. Ленина, 1"},
		{ID: 3, Name: "Склад", Description: "Основной склад", Address: "ул. Складская, 5"},
		{ID: 4, Name: "Офис разработки", Description: "Отдел разработки", Address: "ул. Ленина, 1"},
		{ID: 5, Name: "Бухгалтерия", Description: "Бухгалтерия и финансовый отдел", Address: "ул. Ленина, 1"},
	}
}

// GetSuppliers возвращает список поставщиков
func GetSuppliers() []model.Supplier {
	return []model.Supplier{
		{ID: 1, Name: "ООО Техника", Description: "Поставщик компьютерной техники", Address: "ул. Поставщиков, 1", Phone: "+7 (999) 123-45-67"},
		{ID: 2, Name: "ИП Мебель", Description: "Поставщик офисной мебели", Address: "ул. Мебельная, 10", Phone: "+7 (999) 234-56-78"},
		{ID: 3, Name: "АО Сервер", Description: "Поставщик серверного оборудования", Address: "ул. Серверная, 15", Phone: "+7 (999) 345-67-89"},
		{ID: 4, Name: "ООО Безопасность", Description: "Поставщик систем безопасности", Address: "ул. Охранная, 20", Phone: "+7 (999) 456-78-90"},
		{ID: 5, Name: "ЗАО Сеть", Description: "Поставщик сетевого оборудования", Address: "ул. Сетевая, 25", Phone: "+7 (999) 567-89-01"},
	}
}

// GetEquipment возвращает список оборудования
func GetEquipment() []model.Equipment {
	return []model.Equipment{
		// Компьютерная техника
		{
			Name:         "Ноутбук Dell Latitude 5520",
			Description:  "15.6\" FHD, Intel Core i5-1135G7, 16GB RAM, 512GB SSD",
			SerialNumber: "DELL-2023-001",
			Status:       "available",
			Quantity:     5,
			Price:        89999.99,
			CategoryID:   1,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Компьютер HP ProDesk 600 G6",
			Description:  "Intel Core i7-10700, 32GB RAM, 1TB SSD, Windows 10 Pro",
			SerialNumber: "HP-2023-001",
			Status:       "in_use",
			Quantity:     10,
			Price:        79999.99,
			CategoryID:   1,
			LocationID:   4,
			SupplierID:   1,
		},
		{
			Name:         "Монитор Dell P2419H",
			Description:  "24\" IPS, 1920x1080, 60Hz",
			SerialNumber: "MON-2023-001",
			Status:       "available",
			Quantity:     15,
			Price:        19999.99,
			CategoryID:   1,
			LocationID:   4,
			SupplierID:   1,
		},
		// Сетевое оборудование
		{
			Name:         "Cisco Catalyst 2960-X",
			Description:  "48 портов, управляемый коммутатор",
			SerialNumber: "CISCO-2023-001",
			Status:       "in_use",
			Quantity:     3,
			Price:        149999.99,
			CategoryID:   2,
			LocationID:   2,
			SupplierID:   5,
		},
		{
			Name:         "Ubiquiti UniFi AP AC Pro",
			Description:  "Точка доступа Wi-Fi 5",
			SerialNumber: "UBNT-2023-001",
			Status:       "available",
			Quantity:     8,
			Price:        15999.99,
			CategoryID:   2,
			LocationID:   1,
			SupplierID:   5,
		},
		// Офисная техника
		{
			Name:         "МФУ HP LaserJet Pro M428fdn",
			Description:  "Лазерное МФУ, A4, сетевой",
			SerialNumber: "HPPR-2023-001",
			Status:       "in_use",
			Quantity:     4,
			Price:        39999.99,
			CategoryID:   3,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Сканер Epson WorkForce DS-870",
			Description:  "Документ-сканер, A4, двусторонний",
			SerialNumber: "EPSN-2023-001",
			Status:       "available",
			Quantity:     2,
			Price:        89999.99,
			CategoryID:   3,
			LocationID:   5,
			SupplierID:   1,
		},
		// Мебель
		{
			Name:         "Кресло Ergohuman Plus",
			Description:  "Эргономичное офисное кресло",
			SerialNumber: "CHAIR-2023-001",
			Status:       "in_use",
			Quantity:     20,
			Price:        29999.99,
			CategoryID:   4,
			LocationID:   4,
			SupplierID:   2,
		},
		{
			Name:         "Стол с регулируемой высотой",
			Description:  "Электрический стол-трансформер",
			SerialNumber: "DESK-2023-001",
			Status:       "available",
			Quantity:     15,
			Price:        49999.99,
			CategoryID:   4,
			LocationID:   4,
			SupplierID:   2,
		},
		// Инструменты
		{
			Name:         "Набор инструментов iFixit Pro Tech",
			Description:  "Профессиональный набор для ремонта электроники",
			SerialNumber: "TOOL-2023-001",
			Status:       "in_use",
			Quantity:     3,
			Price:        19999.99,
			CategoryID:   5,
			LocationID:   2,
			SupplierID:   1,
		},
		// Программное обеспечение
		{
			Name:         "Microsoft Office 365 Business",
			Description:  "Годовая подписка на пакет Office 365",
			SerialNumber: "MS365-2023-001",
			Status:       "in_use",
			Quantity:     50,
			Price:        5999.99,
			CategoryID:   6,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Adobe Creative Cloud",
			Description:  "Годовая подписка на пакет Adobe CC",
			SerialNumber: "ADOBE-2023-001",
			Status:       "in_use",
			Quantity:     10,
			Price:        39999.99,
			CategoryID:   6,
			LocationID:   4,
			SupplierID:   1,
		},
		// Телекоммуникационное оборудование
		{
			Name:         "Polycom Studio",
			Description:  "Система видеоконференцсвязи",
			SerialNumber: "POLY-2023-001",
			Status:       "available",
			Quantity:     3,
			Price:        99999.99,
			CategoryID:   7,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Yealink T58A",
			Description:  "IP-телефон с Android",
			SerialNumber: "YEAL-2023-001",
			Status:       "in_use",
			Quantity:     25,
			Price:        19999.99,
			CategoryID:   7,
			LocationID:   1,
			SupplierID:   5,
		},
		// Системы безопасности
		{
			Name:         "Hikvision DS-2CD2143G2-I",
			Description:  "IP-камера 4MP с ИК-подсветкой",
			SerialNumber: "HIK-2023-001",
			Status:       "in_use",
			Quantity:     10,
			Price:        15999.99,
			CategoryID:   8,
			LocationID:   1,
			SupplierID:   4,
		},
		{
			Name:         "СКУД PERCo-Web",
			Description:  "Система контроля доступа",
			SerialNumber: "PERCO-2023-001",
			Status:       "in_use",
			Quantity:     1,
			Price:        299999.99,
			CategoryID:   8,
			LocationID:   1,
			SupplierID:   4,
		},
		// Дополнительное оборудование
		{
			Name:         "ИБП APC Smart-UPS 3000VA",
			Description:  "Источник бесперебойного питания",
			SerialNumber: "APC-2023-001",
			Status:       "in_use",
			Quantity:     4,
			Price:        89999.99,
			CategoryID:   1,
			LocationID:   2,
			SupplierID:   1,
		},
		{
			Name:         "Проектор Epson EB-2250U",
			Description:  "WUXGA проектор 5000 ANSI",
			SerialNumber: "PROJ-2023-001",
			Status:       "available",
			Quantity:     2,
			Price:        129999.99,
			CategoryID:   1,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Сетевое хранилище Synology DS1821+",
			Description:  "NAS на 8 дисков",
			SerialNumber: "SYN-2023-001",
			Status:       "in_use",
			Quantity:     1,
			Price:        159999.99,
			CategoryID:   1,
			LocationID:   2,
			SupplierID:   1,
		},
		{
			Name:         "Серверная стойка 42U",
			Description:  "Стойка для серверного оборудования",
			SerialNumber: "RACK-2023-001",
			Status:       "in_use",
			Quantity:     2,
			Price:        49999.99,
			CategoryID:   4,
			LocationID:   2,
			SupplierID:   2,
		},
		{
			Name:         "Система кондиционирования серверной",
			Description:  "Прецизионный кондиционер",
			SerialNumber: "COOL-2023-001",
			Status:       "in_use",
			Quantity:     2,
			Price:        399999.99,
			CategoryID:   1,
			LocationID:   2,
			SupplierID:   1,
		},
		{
			Name:         "Шредер Fellowes 225Ci",
			Description:  "Офисный шредер, 4 уровень секретности",
			SerialNumber: "SHRED-2023-001",
			Status:       "in_use",
			Quantity:     3,
			Price:        59999.99,
			CategoryID:   3,
			LocationID:   5,
			SupplierID:   1,
		},
		{
			Name:         "Мобильная стойка для ВКС",
			Description:  "Стойка для системы видеоконференцсвязи",
			SerialNumber: "VCST-2023-001",
			Status:       "available",
			Quantity:     2,
			Price:        29999.99,
			CategoryID:   4,
			LocationID:   1,
			SupplierID:   2,
		},
		{
			Name:         "Планшет Samsung Galaxy Tab S7",
			Description:  "Android планшет для презентаций",
			SerialNumber: "TAB-2023-001",
			Status:       "in_use",
			Quantity:     5,
			Price:        49999.99,
			CategoryID:   1,
			LocationID:   1,
			SupplierID:   1,
		},
		{
			Name:         "Система хранения инструментов",
			Description:  "Металлический шкаф с ящиками",
			SerialNumber: "TOOLS-2023-001",
			Status:       "in_use",
			Quantity:     2,
			Price:        39999.99,
			CategoryID:   4,
			LocationID:   2,
			SupplierID:   2,
		},
		{
			Name:         "Тестер кабельных сетей Fluke Networks",
			Description:  "Профессиональный тестер для сетей",
			SerialNumber: "TEST-2023-001",
			Status:       "available",
			Quantity:     1,
			Price:        199999.99,
			CategoryID:   5,
			LocationID:   2,
			SupplierID:   5,
		},
		{
			Name:         "Система видеонаблюдения",
			Description:  "Комплект видеонаблюдения на 16 камер",
			SerialNumber: "CCTV-2023-001",
			Status:       "in_use",
			Quantity:     1,
			Price:        299999.99,
			CategoryID:   8,
			LocationID:   1,
			SupplierID:   4,
		},
		{
			Name:         "Конференц-система Bosch",
			Description:  "Система для конференц-зала",
			SerialNumber: "CONF-2023-001",
			Status:       "in_use",
			Quantity:     1,
			Price:        499999.99,
			CategoryID:   7,
			LocationID:   1,
			SupplierID:   1,
		},
	}
}

// SeedDatabase заполняет базу данных тестовыми данными
func SeedDatabase(db *gorm.DB) error {
	// Создаем категории
	categories := GetCategories()
	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	// Создаем местоположения
	locations := GetLocations()
	if err := db.Create(&locations).Error; err != nil {
		return err
	}

	// Создаем поставщиков
	suppliers := GetSuppliers()
	if err := db.Create(&suppliers).Error; err != nil {
		return err
	}

	// Создаем оборудование
	equipment := GetEquipment()
	for i := range equipment {
		equipment[i].CreatedAt = time.Now()
		equipment[i].UpdatedAt = time.Now()
	}
	if err := db.Create(&equipment).Error; err != nil {
		return err
	}

	return nil
}
