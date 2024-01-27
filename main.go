package main

import (
	".webapp/middlewares"
	".webapp/models"
	".webapp/routes"
	"encoding/csv"
	"fmt"
	"github.com/alexcesaro/statsd"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"os"
	"time"
)

var r *gin.Engine

func main() {
	//f, err := os.OpenFile("csye6225.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	f, err := os.OpenFile("/var/log/webapp/csye6225.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// 设置gin日志文件
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 设置默认log的文件
	log.SetOutput(f)
	gin.SetMode(gin.ReleaseMode)

	err = models.Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("error boot db: %v", err)
	}
	err = models.Db.AutoMigrate(&models.Assignment{})
	if err != nil {
		log.Fatalf("error boot db: %v", err)
	}
	err = models.Db.AutoMigrate(&models.UserAssignment{})
	if err != nil {
		log.Fatalf("error boot db: %v", err)
	}
	err = models.Db.AutoMigrate(&models.Submission{})
	if err != nil {
		log.Fatalf("error boot db: %v", err)
	}
	readUsers()
	r = gin.Default()

	// 创建StatsD客户端
	statsdClient, err := statsd.New(
		statsd.Address("127.0.0.1:8125"), // 替换为您的StatsD服务器地址
		statsd.Prefix("api"),             // 统计数据的前缀
	)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer statsdClient.Close()

	// 使用中间件记录每个API请求的次数
	r.Use(func(c *gin.Context) {

		// 使用请求路径作为统计的唯一标识
		method := c.Request.Method
		path := c.Request.URL.Path

		// 创建带有HTTP方法和路径的度量名称
		metricName := fmt.Sprintf("request.count.%s%s", method, path)

		// 统计每个API请求次数，包含HTTP方法
		statsdClient.Increment(metricName)

		// 继续处理请求
		c.Next()
	})
	r.Use(middlewares.CheckMethodNotAllowed(r))
	routes.AssignmentRoutesInit(r)
	routes.HealthCheckInit(r)
	log.Println("API server starting...")
	r.Run("0.0.0.0:8080")
}

func readUsers() {
	file, err := os.Open("/opt/csye6225/users.csv")
	//file, err := os.Open("users.csv")
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading CSV row:", err)
		}
		if lineNumber == 0 {
			lineNumber++
			continue
		}
		var existingUser models.User
		result := models.Db.Where("email = ?", row[2]).First(&existingUser)
		if result.Error == nil {
			continue
		}
		user := models.User{
			FirstName: row[0],
			LastName:  row[1],
			Email:     row[2],
		}
		password, pswdErr := HashPassword(row[3])
		if pswdErr != nil {
			log.Fatal("Error hashing password:", pswdErr)
		}
		user.Password = password
		user.AccountCreated = time.Now()
		user.AccountUpdated = time.Now()
		result = models.Db.Create(&user)
		if result.Error != nil {
			log.Printf("Create user failed: %v", result.Error)
		}
	}
}

func HashPassword(password string) (string, error) {
	// 生成密码的 Bcrypt 哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// 返回密码
	return string(hashedPassword), nil
}
