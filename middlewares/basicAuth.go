package middlewares

import (
	".webapp/models"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func InitBasicAuth(c *gin.Context) {
	// 从请求头中获取 Authorization 字段
	authHeader := c.Request.Header.Get("Authorization")

	// 检查是否提供了基本认证信息
	if authHeader == "" {
		// 如果没有提供认证信息，返回未授权响应
		c.Header("WWW-Authenticate", "Basic realm=Restricted")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 解析基本认证信息
	email, password, ok := parseBasicAuth(authHeader)
	if !ok {
		// 如果无法解析认证信息，返回未授权响应
		c.Header("WWW-Authenticate", "Basic realm=Restricted")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 在此处验证用户名和密码，你可以将它们与存储的用户凭据进行比较
	// 如果验证成功，允许请求继续；否则，返回未授权响应
	valid, id := isValidUser(email, password)
	if !valid {
		c.Header("WWW-Authenticate", "Basic realm=Restricted")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userID", id)
	// 验证通过，继续执行下一个处理器
	c.Next()
}

// isValidUser 验证用户名和密码是否有效
func isValidUser(email, password string) (bool, string) {
	// 从数据库中检索与提供的用户名匹配的用户记录，以获取存储的哈希密码
	user, err := getUserByEmail(email)
	if err != nil {
		return false, "" // 用户不存在或发生错误
	}

	// 使用 bcrypt.CompareHashAndPassword 检查密码是否匹配
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, ""
	}
	return true, user.ID
}

// getUserByEmail 根据用户名从数据库中检索用户记录
func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := models.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func parseBasicAuth(authHeader string) (string, string, bool) {
	// 检查 Authorization 字段是否以 "Basic " 开头
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", false
	}

	// 提取 base64 编码的用户名和密码部分
	encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
	// 解码 base64 字符串
	decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return "", "", false
	}

	// 将解码后的字符串拆分为用户名和密码
	creds := strings.SplitN(string(decodedCreds), ":", 2)
	if len(creds) != 2 {
		return "", "", false
	}

	// 返回用户名、密码和解析成功标志
	return creds[0], creds[1], true
}
