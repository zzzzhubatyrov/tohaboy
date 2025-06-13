package service

import (
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
	newUser := &model.User{
		Username: user["username"],
		Password: user["password"],
	}

	login, err := s.repo.Login(newUser)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(user["password"])); err != nil {
		log.Println("[service] invalid password:", err)
		return nil, err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(login.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	token, err := claims.SignedString([]byte(viper.GetString("SECRET_KEY")))
	if err != nil {
		log.Println("[service] could not generate token:", err)
		return nil, err
	}

	return &model.LoginResponse{
		User:  login,
		Token: token,
	}, nil
}

func (s *AuthService) Register(user map[string]string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user["password"]), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Username: user["username"],
		Password: string(hashedPassword),
	}

	reg, err := s.repo.Register(newUser)
	if err != nil {
		return nil, err
	}

	return reg, nil
}
