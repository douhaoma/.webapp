package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 如果路由没有匹配到任何处理程序并且请求方法是不允许的方法，则返回 405
		if c.Writer.Status() == 404 {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method Not Allowed"})
			c.Abort()
		}
	}
}

func CheckMethodNotAllowed(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的路径
		path := c.Request.URL.Path

		// 获取所有注册的路由
		routes := r.Routes()

		var methodExists bool
		// 检查是否存在相同路径的路由
		for _, route := range routes {
			if path == route.Path {
				methodExists = true
				// 如果方法匹配，则直接继续处理请求
				if c.Request.Method == route.Method {
					c.Next()
					return
				}
			}
		}

		// 如果存在相同路径但方法不匹配，则返回405错误
		if methodExists {
			c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
			c.Abort()
			return
		}

		// 如果没有找到匹配的路径，则继续处理请求
		c.Next()
	}
}

//func InitMethodNotAllowed(c *gin.Context) {
//	c.Next()
//	// 如果路由没有匹配到任何处理程序并且请求方法是不允许的方法，则返回 405
//	if c.Writer.Status() == 404 {
//		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method Not Allowed"})
//		c.Abort()
//	}
//}
