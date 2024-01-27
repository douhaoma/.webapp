package controllers

import (
	".webapp/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("X-Content-Type-Options", "nosniff")
	l := len(c.Request.URL.Query())
	if l > 0 {
		// 如果请求体不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query not allowed for this request"})
		return
	}
	contentLength := c.Request.ContentLength
	if contentLength > 0 {
		// 如果请求体不为空，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body not allowed for GET request"})
		return
	}
	sqlDB, err := models.Db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{})
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
