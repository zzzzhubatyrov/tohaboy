package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Date специальный тип для работы с датами
type Date time.Time

// MarshalJSON реализует интерфейс json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
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

	*d = Date(parsedTime)
	return nil
}

// Scan реализует интерфейс sql.Scanner
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date(v)
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*d = Date(parsedTime)
	}
	return nil
}

// Value реализует интерфейс driver.Valuer
func (d Date) Value() (interface{}, error) {
	return time.Time(d), nil
}

// Time конвертирует Date в time.Time
func (d Date) Time() time.Time {
	return time.Time(d)
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
	ID        int            `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique" json:"username"`
	Password  string         `json:"password"`
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
	ID           int            `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	SerialNumber string         `gorm:"unique" json:"serial_number"`
	Category     string         `json:"category"`
	Description  string         `json:"description"`
	Price        float64        `json:"price"`
	Status       string         `json:"status"`
	Quantity     int            `json:"quantity"`
	LocationID   int            `json:"location_id"`
	Location     *Location      `gorm:"foreignKey:LocationID" json:"location"`
	SupplierID   int            `json:"supplier_id"`
	Supplier     *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier"`
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
//	ContactInfo - контакты (телефон/email/контактное лицо)
//	Equipment - список поставляемого оборудования
type Supplier struct {
	ID          int         `gorm:"primaryKey" json:"id"`
	Name        string      `json:"name"`
	ContactInfo string      `json:"contact_info"`
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
	ID            int         `gorm:"primaryKey" json:"id"`
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
	ID             int        `gorm:"primaryKey" json:"id"`
	EquipmentID    int        `json:"equipment_id"`
	Equipment      *Equipment `gorm:"foreignKey:EquipmentID;references:ID" json:"equipment"`
	FromLocationID int        `json:"from_location_id"`
	FromLocation   *Location  `gorm:"foreignKey:FromLocationID;references:ID" json:"from_location"`
	ToLocationID   int        `json:"to_location_id"`
	ToLocation     *Location  `gorm:"foreignKey:ToLocationID;references:ID" json:"to_location"`
	Quantity       int        `json:"quantity"`
	Reason         string     `json:"reason"`
	CreatedByID    int        `json:"created_by_id"`
	CreatedBy      *User      `gorm:"foreignKey:CreatedByID;references:ID" json:"created_by"`
	DocumentID     int        `json:"document_id"`
	Date           Date       `json:"date" gorm:"type:date"`
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
	ID           int            `gorm:"primaryKey" json:"id"`
	Type         string         `json:"type"`
	Number       string         `gorm:"unique" json:"number"`
	Status       string         `json:"status"`
	Date         Date           `json:"date" gorm:"type:date"`
	CreatedByID  int            `json:"created_by_id"`
	CreatedBy    *User          `gorm:"foreignKey:CreatedByID" json:"created_by"`
	ApprovedByID int            `gorm:"default:null" json:"approved_by_id"`
	ApprovedBy   *User          `gorm:"foreignKey:ApprovedByID" json:"approved_by"`
	LocationID   int            `json:"location_id"`
	Location     *Location      `gorm:"foreignKey:LocationID" json:"location"`
	Items        []DocumentItem `gorm:"foreignKey:DocumentID" json:"items"`
	Comment      string         `json:"comment"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// DocumentItem представляет элемент документа (приход/расход/инвентаризация)
//
//	ID - уникальный идентификатор
//	DocumentID - ссылка на документ
//	EquipmentID - ссылка на оборудование
//	Equipment - связанное оборудование (gorm relation)
//	Quantity - количество для документов прихода/расхода
//	ActualQuantity - фактическое количество при инвентаризации
//	Price - цена за единицу
//	TotalPrice - общая сумма
//	Comment - комментарий
type DocumentItem struct {
	ID             int        `gorm:"primaryKey" json:"id"`
	DocumentID     int        `json:"document_id"`
	EquipmentID    int        `json:"equipment_id"`
	Equipment      *Equipment `gorm:"foreignKey:EquipmentID" json:"equipment"`
	Quantity       int        `json:"quantity"`        // Количество для документов прихода/расхода
	ActualQuantity int        `json:"actual_quantity"` // Фактическое количество при инвентаризации
	Price          float64    `json:"price"`           // Цена за единицу
	TotalPrice     float64    `json:"total_price"`     // Общая сумма
	Comment        string     `json:"comment"`
}
