package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Date специальный тип для работы с датами
type Date struct {
	t time.Time
}

// NewDate создает новый объект Date
func NewDate(t time.Time) Date {
	return Date{t: t}
}

// MarshalJSON реализует интерфейс json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.t.Format("2006-01-02"))
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	d.t = parsedTime
	return nil
}

// Scan реализует интерфейс sql.Scanner
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.t = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.t = v
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.t = parsedTime
	}
	return nil
}

// Value реализует интерфейс driver.Valuer
func (d Date) Value() (interface{}, error) {
	return d.t, nil
}

// Time возвращает time.Time
func (d Date) Time() time.Time {
	return d.t
}

// User представляет учётную запись пользователя системы
// Поля:
//
//	ID - уникальный идентификатор пользователя
//	Username - логин пользователя (уникальный)
//	Password - хэш пароля (не должен возвращаться в JSON)
//	Role - роль пользователя: "admin", "manager", "auditor"
//	CreatedAt - дата создания учетной записи
//	UpdatedAt - дата последнего обновления
//	DeletedAt - метка мягкого удаления (не экспортируется в JSON)
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique" json:"username"`
	Password  string         `json:"password"` // Don't expose password in JSON
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Equipment описывает единицу оборудования/материала
// Поля:
//
//	ID - уникальный идентификатор
//	Name - название оборудования
//	SerialNumber - серийный номер (уникальный)
//	Category - категория оборудования
//	Description - описание/характеристики
//	Price - стоимость единицы (в валюте)
//	Status - текущий статус: "available", "in_use", "maintenance", "written_off"
//	Quantity - общее количество на складе
//	LocationID - ссылка на местоположение
//	Location - связанное местоположение (gorm relation)
//	SupplierID - ссылка на поставщика
//	Supplier - связанный поставщик (gorm relation)
//	Movements - история перемещений
//	Documents - связанные документы
//	CreatedAt/UpdatedAt - метки времени
type Equipment struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	SerialNumber string         `gorm:"unique" json:"serial_number"`
	Status       string         `json:"status"`
	Quantity     int            `json:"quantity"`
	Price        float64        `json:"price"`
	CategoryID   uint           `json:"category_id"`
	Category     *Category      `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	LocationID   uint           `json:"location_id"`
	Location     *Location      `gorm:"foreignKey:LocationID;references:ID" json:"location"`
	SupplierID   uint           `json:"supplier_id"`
	Supplier     *Supplier      `gorm:"foreignKey:SupplierID;references:ID" json:"supplier"`
	Movements    []Movement     `gorm:"foreignKey:EquipmentID" json:"movements"`
	Documents    []DocumentItem `gorm:"foreignKey:EquipmentID" json:"documents"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Supplier содержит информацию о поставщике оборудования
// Поля:
//
//	ID - уникальный идентификатор
//	Name - название компании
//	Description - описание компании
//	Address - физический адрес компании
//	Phone - контактный телефон компании
//	Equipment - список поставляемого оборудования
type Supplier struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Address     string      `json:"address"`
	Phone       string      `json:"phone"`
	Equipment   []Equipment `gorm:"foreignKey:SupplierID" json:"equipment"`
}

// Location описывает место хранения оборудования
// Поля:
//
//	ID - уникальный идентификатор
//	Name - название места (склад/кабинет)
//	Description - описание места
//	Address - физический адрес
//	Equipment - список оборудования на этом месте
//	FromMovements - история перемещений из этого места
//	ToMovements - история перемещений в это место
type Location struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Address       string      `json:"address"`
	Equipment     []Equipment `gorm:"foreignKey:LocationID" json:"equipment"`
	FromMovements []Movement  `gorm:"foreignKey:FromLocationID" json:"from_movements"`
	ToMovements   []Movement  `gorm:"foreignKey:ToLocationID" json:"to_movements"`
}

// Movement фиксирует факт перемещения оборудования
// Поля:
//
//	ID - уникальный идентификатор
//	EquipmentID - ссылка на оборудование
//	FromLocationID - откуда перемещается (0 если приемка)
//	ToLocationID - куда перемещается
//	Quantity - количество перемещаемых единиц
//	Reason - причина: "transfer", "inventory", "repair"
//	CreatedByID - кто создал перемещение
//	CreatedBy - связанный пользователь (создатель)
//	DocumentID - ссылка на документ-основание
//	Date - дата перемещения
type Movement struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	EquipmentID    uint       `json:"equipment_id"`
	Equipment      *Equipment `gorm:"foreignKey:EquipmentID;references:ID" json:"equipment"`
	FromLocationID uint       `json:"from_location_id"`
	FromLocation   *Location  `gorm:"foreignKey:FromLocationID;references:ID" json:"from_location"`
	ToLocationID   uint       `json:"to_location_id"`
	ToLocation     *Location  `gorm:"foreignKey:ToLocationID;references:ID" json:"to_location"`
	Quantity       int        `json:"quantity"`
	Reason         string     `json:"reason"`
	CreatedByID    uint       `json:"created_by_id"`
	CreatedBy      *User      `gorm:"foreignKey:CreatedByID;references:ID" json:"created_by"`
	DocumentID     uint       `json:"document_id"`
	Date           time.Time  `json:"date" gorm:"type:date"`
}

// MarshalJSON реализует интерфейс json.Marshaler для Movement
func (m *Movement) MarshalJSON() ([]byte, error) {
	type Alias Movement
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(m),
		Date:  m.Date.Format("2006-01-02"),
	})
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler для Movement
func (m *Movement) UnmarshalJSON(data []byte) error {
	type Alias Movement
	aux := &struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Date != "" {
		parsedTime, err := time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return err
		}
		m.Date = parsedTime
	}
	return nil
}

// Document представляет документ учета (акт/накладная)
// Поля:
//
//	ID - уникальный идентификатор
//	Type - тип документа: "inventory", "transfer", "write_off", "acceptance"
//	Number - уникальный номер документа (формат "ИНВ-2023-001")
//	Status - статус: "draft", "completed", "canceled"
//	Date - дата документа
//	CreatedByID - кто создал документ
//	CreatedBy - связанный пользователь (создатель)
//	ApprovedByID - кто утвердил (может быть null)
//	ApprovedBy - связанный пользователь (утвердивший)
//	LocationID - место проведения (для инвентаризации)
//	Location - связанное местоположение
//	Items - позиции документа
//	Comment - комментарий к документу
//	CreatedAt/UpdatedAt - метки времени
type Document struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Type         string         `json:"type"`
	Number       string         `gorm:"unique" json:"number"`
	Status       string         `json:"status"`
	Date         time.Time      `json:"date" gorm:"type:date"`
	CreatedByID  uint           `json:"created_by_id"`
	CreatedBy    *User          `gorm:"foreignKey:CreatedByID" json:"created_by"`
	ApprovedByID uint           `gorm:"default:null" json:"approved_by_id"`
	ApprovedBy   *User          `gorm:"foreignKey:ApprovedByID" json:"approved_by"`
	LocationID   uint           `json:"location_id"`
	Location     *Location      `gorm:"foreignKey:LocationID" json:"location"`
	Items        []DocumentItem `gorm:"foreignKey:DocumentID" json:"items"`
	Comment      string         `json:"comment"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// MarshalJSON реализует интерфейс json.Marshaler для Document
func (d *Document) MarshalJSON() ([]byte, error) {
	type Alias Document
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(d),
		Date:  d.Date.Format("2006-01-02"),
	})
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler для Document
func (d *Document) UnmarshalJSON(data []byte) error {
	type Alias Document
	aux := &struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Date != "" {
		parsedTime, err := time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return err
		}
		d.Date = parsedTime
	}
	return nil
}

// DocumentItem представляет позицию документа
// Поля:
//
//	ID - уникальный идентификатор
//	DocumentID - ссылка на документ
//	EquipmentID - ссылка на оборудование
//	Equipment - связанное оборудование
//	Quantity - количество
//	ActualQuantity - фактическое количество
//	Price - цена единицы
//	TotalPrice - общая стоимость
//	Comment - комментарий
type DocumentItem struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID     uint      `gorm:"not null;index:idx_doc_equipment,priority:1" json:"document_id"`
	EquipmentID    uint      `gorm:"not null;index:idx_doc_equipment,priority:2" json:"equipment_id"`
	Equipment      Equipment `gorm:"foreignKey:EquipmentID;references:ID" json:"equipment"`
	Quantity       int       `gorm:"not null" json:"quantity"`
	ActualQuantity int       `json:"actual_quantity"`
	Price          float64   `gorm:"not null" json:"price"`
	TotalPrice     float64   `gorm:"not null" json:"total_price"`
	Comment        string    `json:"comment"`
}

// Category represents an equipment category
type Category struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryResponse represents a response containing a single category
type CategoryResponse struct {
	Model   *Category `json:"model"`
	Message string    `json:"msg"`
}

// CategoryListResponse represents a response containing a list of categories
type CategoryListResponse struct {
	Model   []Category `json:"model"`
	Message string     `json:"msg"`
}
