package service

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"tohaboy/internal/model"
	"tohaboy/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.AuthRepositoryInterface
}

func NewAuthService(repo repository.AuthRepositoryInterface) *AuthService {
	service := &AuthService{repo: repo}
	// Создаем администратора при инициализации сервиса
	if err := service.createAdminIfNotExists(); err != nil {
		log.Printf("[service] error creating admin: %v", err)
	}
	return service
}

func (s *AuthService) createAdminIfNotExists() error {
	// Проверяем, существует ли уже пользователь admin
	adminUser := &model.User{
		Username: "admin",
	}

	existingAdmin, err := s.repo.Login(adminUser)
	if err == nil && existingAdmin != nil {
		// Администратор уже существует
		return nil
	}

	// Создаем нового администратора
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newAdmin := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	_, err = s.repo.Register(newAdmin)
	if err != nil {
		return err
	}

	log.Println("[service] admin user created successfully")
	return nil
}

func (s *AuthService) Login(user map[string]string) (*model.LoginResponse, error) {
	// Validate required fields
	if user["username"] == "" || user["password"] == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	newUser := &model.User{
		Username: user["username"],
	}

	login, err := s.repo.Login(newUser)
	if err != nil {
		log.Printf("[service] login error: %v", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(user["password"])); err != nil {
		log.Printf("[service] invalid password: %v", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.FormatUint(uint64(login.ID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	token, err := claims.SignedString([]byte(viper.GetString("SECRET_KEY")))
	if err != nil {
		log.Printf("[service] could not generate token: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	// Don't expose password hash in response
	login.Password = ""
	return &model.LoginResponse{
		User:  login,
		Token: token,
	}, nil
}

func (s *AuthService) Register(user map[string]string) (*model.User, error) {
	// Validate required fields
	if user["username"] == "" || user["password"] == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user["password"]), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	newUser := &model.User{
		Username:  user["username"],
		Password:  string(hashedPassword),
		Role:      user["role"], // If role is not provided, it will be empty
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	reg, err := s.repo.Register(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	// Don't return the password hash
	reg.Password = ""
	return reg, nil
}
