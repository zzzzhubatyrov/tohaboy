package model

// LoginResponse содержит данные для авторизованного пользователя
// Поля:
//
//	User - объект пользователя (без sensitive-полей)
//	Token - JWT-токен для аутентификации
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// Response базовый тип ответа
type Response[T any] struct {
	Model   T      `json:"model"`
	Message string `json:"msg"`
}

// Конкретные типы для каждого вида ответа
type EquipmentResponse struct {
	Model   *Equipment `json:"model"`
	Message string     `json:"msg"`
}

type EquipmentListResponse struct {
	Model   []Equipment `json:"model"`
	Message string      `json:"msg"`
}

type DocumentResponse struct {
	Model   *Document `json:"model"`
	Message string    `json:"msg"`
}

type DocumentListResponse struct {
	Model   []Document `json:"model"`
	Message string     `json:"msg"`
}

type MovementResponse struct {
	Model   *Movement `json:"model"`
	Message string    `json:"msg"`
}

type MovementListResponse struct {
	Model   []Movement `json:"model"`
	Message string     `json:"msg"`
}

type SupplierResponse struct {
	Model   *Supplier `json:"model"`
	Message string    `json:"msg"`
}

type SupplierListResponse struct {
	Model   []Supplier `json:"model"`
	Message string     `json:"msg"`
}

type LocationResponse struct {
	Model   *Location `json:"model"`
	Message string    `json:"msg"`
}

type LocationListResponse struct {
	Model   []Location `json:"model"`
	Message string     `json:"msg"`
}
