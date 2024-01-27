package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Assignment struct {
	ID                 string    `json:"id" gorm:"primary_key"`
	Name               string    `json:"name" gorm:"not null"`
	Points             uint      `gorm:"check:points >= 1 AND points <= 100;not null" json:"points"`
	Num_of_attempts    uint      `gorm:"check:num_of_attempts >= 1 AND num_of_attempts <= 100;not null" json:"num_of_attempts"`
	Deadline           time.Time `json:"deadline" gorm:"not null"`
	Assignment_created time.Time `json:"assignment_created" gorm:"autoCreateTime"`
	Assignment_updated time.Time `json:"assignment_updated" gorm:"autoUpdateTime"`
}

// GetAllAssignments 从数据库中获取所有作业的列表
func GetAllAssignments() ([]Assignment, error) {
	var assignments []Assignment
	result := Db.Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
	}
	return assignments, nil
}

func CreateAssignment(c *gin.Context, assignment *Assignment) error {
	userID := c.GetString("userID")
	tx := Db.Begin()
	//fmt.Println(assignment.Assignment_created)
	if err := tx.Create(&assignment).Error; err != nil {
		fmt.Println("1")
		tx.Rollback() // 发生错误时回滚事务
		return err
	}

	userAssignment := UserAssignment{
		UserID:       userID,
		AssignmentID: assignment.ID,
	}
	result := Db.Create(&userAssignment)
	if result.Error != nil {
		fmt.Println(result.Error)
		fmt.Println("2")
		tx.Rollback() // 发生错误时回滚事务
		return result.Error
	}
	tx.Commit()
	return nil
}

func GetAssignmentById(id string) (*Assignment, error) {
	var assignment Assignment
	result := Db.Where("id = ?", id).First(&assignment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assignment, nil
}

func DeleteAssignmentById(id string) error {
	var assignment Assignment
	err := Db.Where("id = ?", id).Delete(&assignment).Error
	return err
}
