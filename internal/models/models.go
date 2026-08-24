package models

import (
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
}

type Scenario struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Kind      string    `gorm:"not null"` // fire|sip|emi|emergency|budget
	Title     string    `gorm:"not null"`
	Inputs    string    `gorm:"not null"`
	Outputs   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func Open(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &Scenario{}); err != nil {
		return nil, err
	}
	return db, nil
}
