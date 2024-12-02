package models

import (
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            int    `gorm:"primaryKey"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required" gorm:"unique"`
	Password      string `json:"password" binding:"required"`
	WorkspaceName string `json:"workspaceName,omitempty"`
}

func (user *User) CreateUserRecord() error {
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
