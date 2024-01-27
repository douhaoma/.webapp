package models

import (
	"github.com/gin-gonic/gin"
	"time"
)

type User struct {
	ID             string `gorm:"primaryKey"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Password       string `gorm:"not null"`
	Email          string `gorm:"unique; not null"`
	AccountCreated time.Time
	AccountUpdated time.Time
}

func GetUserById(c *gin.Context, id string) (*User, error) {
	var user User
	result := Db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
