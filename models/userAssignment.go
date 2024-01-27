package models

type UserAssignment struct {
	AssignmentID string `gorm:"primaryKey"`
	UserID       string `gorm:"primaryKey"`
}
