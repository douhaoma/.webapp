package models

import (
	"fmt"
	"time"
)

type Submission struct {
	ID                string `json:"id" gorm:"primary_key"`
	UserId            string
	AssignmentId      string    `json:"assignment_id"`
	SubmissionUrl     string    `json:"submission_url"`
	SubmissionDate    time.Time `json:"submission_date" gorm:"autoUpdateTime"`
	SubmissionUpdated time.Time `json:"submission_updated" gorm:"autoUpdateTime"`
}

func (s Submission) GetCurrAttempts(userId string) int {
	var count int64
	Db.Model(&Submission{}).Where("user_id = ? AND assignment_id = ?", userId, s.AssignmentId).Count(&count)
	return int(count)
}

func (s *Submission) CreateSubmission() error {
	result := Db.Create(&s)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

func (s *Submission) GetSubmissionsNums() int {
	var count int64
	Db.Model(&Submission{}).Where("assignment_id = ?", s.AssignmentId).Count(&count)
	return int(count)
}
