package models

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var Db *gorm.DB
var TopicArn, Region string

func init() {
	// read app.ini
	path := "/opt/csye6225/application.properties"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "./config/application.properties.sample"
	}
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	hostname := cfg.Section("mysql").Key("hostname").String()
	username := cfg.Section("mysql").Key("username").String()
	password := cfg.Section("mysql").Key("password").String()
	database := cfg.Section("mysql").Key("database").String()
	TopicArn = cfg.Section("sns").Key("topicArn").String()
	Region = cfg.Section("sns").Key("region").String()
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=true&loc=Local", username, password, hostname, database)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true,
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
	}
	var tableNames []string
	Db.Raw("SHOW TABLES").Pluck("Tables_in_your_database", &tableNames)

	//// 删除所有表格
	//for _, tableName := range tableNames {
	//	Db.Migrator().DropTable(tableName)
	//}

	////
	Db.Callback().Create().Before("gorm:before_create").Register("generateUUID", GenerateUUIDCallback)
}

func GenerateUUIDCallback(db *gorm.DB) {
	schema := db.Statement.Schema
	// 根据表的名称判断是否生成 UUID
	if schema.Table == "users" || schema.Table == "assignments" || schema.Table == "submissions" {
		fmt.Println("1")
		uuidString := uuid.New().String()
		db.Statement.SetColumn("ID", uuidString)
	}
}
