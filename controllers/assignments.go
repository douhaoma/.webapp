package controllers

import (
	".webapp/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func GetAssignmentsList(c *gin.Context) {
	l := len(c.Request.URL.Query())
	if l > 0 {
		// 如果查询不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query not allowed for this request"})
		return
	}
	contentLength := c.Request.ContentLength
	if contentLength > 0 {
		// 如果请求体不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body not allowed for GET request"})
		return
	}
	assignments, err := models.GetAllAssignments()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func CreateAssignment(c *gin.Context) {
	//fmt.Println("1")
	var assignmentJSON struct {
		Name            string    `json:"name"`
		Points          uint      `json:"points"`
		Num_of_attempts uint      `json:"num_of_attempts"`
		Deadline        time.Time `json:"deadline"`
	}

	// 绑定 JSON 数据到 assignmentJSON 结构体
	if err := c.BindJSON(&assignmentJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is not valid"})
		return
	}

	// 创建一个 Assignment 结构体并复制 assignmentJSON 的值
	assignment := &models.Assignment{
		Name:            assignmentJSON.Name,
		Points:          assignmentJSON.Points,
		Num_of_attempts: assignmentJSON.Num_of_attempts,
		Deadline:        assignmentJSON.Deadline,
		//Assignment_created: time.Now(),
		//Assignment_updated: time.Now(),
	}
	result := models.CreateAssignment(c, assignment)
	if result != nil {
		fmt.Println("3")
		fmt.Println(result.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	//fmt.Println(assignment.Assignment_created)
	c.JSON(http.StatusCreated, assignment)
}

func GetAssignmentById(c *gin.Context) {
	contentLength := c.Request.ContentLength
	if contentLength > 0 {
		// 如果请求体不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body not allowed for GET request"})
		return
	}
	//id := c.Query("id")

	id := c.Param("id")
	fmt.Println(">>>>>>>>>>>>>>>>>>", id)
	assignment, err := models.GetAssignmentById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, assignment)
}

func DeleteAssignmentById(c *gin.Context) {
	contentLength := c.Request.ContentLength
	if contentLength > 0 {
		// 如果请求体不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body not allowed for DELETE request"})
		return
	}

	//判断是否由当前用户创造
	userID := c.GetString("userID")
	//id := c.Query("id")
	id := c.Param("id")
	var userAssignment models.UserAssignment
	result := models.Db.Where("user_id = ? AND assignment_id = ?", userID, id).First(&userAssignment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果关系表中没有记录，表示该用户没有创建该作业，返回 403 Forbidden
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have permission to delete this assignment"})
			return
		}
		c.JSON(400, gin.H{})
		return
	}
	//fmt.Println("delete")
	sub := &models.Submission{
		AssignmentId: id,
	}
	currSubmission := sub.GetSubmissionsNums()
	fmt.Println(currSubmission)
	if currSubmission > 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment has submissions, can not delete"})
		return
	}
	//fmt.Println(id)
	err := models.DeleteAssignmentById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	var ua *models.UserAssignment
	models.Db.Where("assignment_id = ?", id).Delete(&ua)
	c.JSON(http.StatusNoContent, gin.H{})
}

func UpdateAssignment(c *gin.Context) {
	//id := c.Query("id")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is missing"})
		return
	}

	var assignmentJSON struct {
		Name            string    `json:"name"`
		Points          uint      `json:"points"`
		Num_of_attempts uint      `json:"num_of_attempts"`
		Deadline        time.Time `json:"deadline"`
	}

	// 绑定 JSON 数据到 assignmentJSON 结构体
	if err := c.BindJSON(&assignmentJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is not valid"})
		return
	}
	userID := c.GetString("userID")

	var userAssignment models.UserAssignment
	//检查关系表中是否有该条数据
	result := models.Db.Where("user_id = ? AND assignment_id = ?", userID, id).First(&userAssignment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果关系表中没有记录，表示该用户没有创建该作业，返回 403 Forbidden
			c.JSON(http.StatusForbidden, gin.H{"error": "User does not have permission to delete this assignment"})
			return
		}
		c.JSON(400, gin.H{})
		return
	}
	assignment := &models.Assignment{
		Name:            assignmentJSON.Name,
		Points:          assignmentJSON.Points,
		Num_of_attempts: assignmentJSON.Num_of_attempts,
		Deadline:        assignmentJSON.Deadline,
	}
	assignment.ID = id
	models.Db.Where("id = ?", id).Updates(&assignment)
	c.JSON(204, gin.H{})
}

type submissionUrl struct {
	Url string `json:"submission_url"`
}

type submissionWithoutUserId struct {
	ID                string    `json:"id" gorm:"primary_key"`
	AssignmentId      string    `json:"assignment_id"`
	SubmissionUrl     string    `json:"submission_url"`
	SubmissionDate    time.Time `json:"submission_date" gorm:"autoUpdateTime"`
	SubmissionUpdated time.Time `json:"submission_updated" gorm:"autoUpdateTime"`
}

type SubmissionMessage struct {
	StudentEmail   string `json:"student_email"`
	SubmissionURL  string `json:"submission_url"`
	AssignmentName string `json:"assignment_name"`
}

func SubmitAssignment(c *gin.Context) {
	assignmentId := c.Param("id")
	userId := c.GetString("userID")
	user, _ := models.GetUserById(c, userId)
	assignment, err := models.GetAssignmentById(assignmentId)
	// 没有这个作业
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No this assignment"})
		return
	}
	// body不对
	var subUrl submissionUrl
	if err := c.BindJSON(&subUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json is not valid"})
		return
	}

	submission := &models.Submission{
		AssignmentId:      assignmentId,
		UserId:            userId,
		SubmissionUrl:     subUrl.Url,
		SubmissionDate:    time.Now(),
		SubmissionUpdated: time.Now(),
	}
	// 超过提交时间
	if submission.SubmissionDate.After(assignment.Deadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exceeded the submission deadline"})
		return
	}

	// 超过提交次数
	numOfAtp := int(assignment.Num_of_attempts)
	fmt.Println(numOfAtp)
	if submission.GetCurrAttempts(userId) >= numOfAtp {
		c.JSON(http.StatusForbidden, gin.H{"error": "Exceeded the attempts limit"})
		return
	}

	// 尝试提交
	err = submission.CreateSubmission()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// 成功
	submission1 := &submissionWithoutUserId{
		ID:                submission.ID,
		AssignmentId:      assignmentId,
		SubmissionUrl:     submission.SubmissionUrl,
		SubmissionUpdated: submission.SubmissionUpdated,
		SubmissionDate:    submission.SubmissionDate,
	}
	req, err := postToSNSTopic(models.TopicArn, models.Region, user.Email, submission.SubmissionUrl, assignment.Name)
	if err != nil {
		log.Printf("WARN: post to SNS topic error: %s\n", err)
	} else {
		log.Printf("INFO post message to aws sns, req: %v", req)
	}
	c.JSON(http.StatusOK, submission1)
}

func postToSNSTopic(topicArn, region, email, url, assignmentName string) (*sns.PublishOutput, error) {

	message := &SubmissionMessage{
		StudentEmail:   email,
		SubmissionURL:  url,
		AssignmentName: assignmentName,
	}

	// 将消息结构体转换为JSON字符串
	messageJson, err := json.Marshal(message)
	// 创建 AWS 会话
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// 创建 SNS 服务客户端
	svc := sns.New(sess)

	// 发布到 SNS 主题
	req, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageJson)),
		TopicArn: aws.String(topicArn),
	})

	return req, err
}
