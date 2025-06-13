package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(storageName string) *Storage {
	db, err := gorm.Open(sqlite.Open(storageName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Storage{db: db}
}

func (s *Storage) GetDB() *gorm.DB {
	return s.db
}

func (s *Storage) Migrate(models []interface{}) error {
	migrator := s.db.Migrator()
	for _, model := range models {
		if err := migrator.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) DropTables(models []interface{}) error {
	migrator := s.db.Migrator()
	for _, model := range models {
		if err := migrator.DropTable(model); err != nil {
			return err
		}
	}
	return nil
}
