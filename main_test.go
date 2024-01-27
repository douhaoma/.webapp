package main

import (
	".webapp/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzEndpoint(t *testing.T) {
	// 创建一个新的 HTTP 请求对象
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个响应记录器
	rr := httptest.NewRecorder()

	// 获取你的应用程序的路由对象
	r := gin.Default()        // 使用 gin.Default() 创建一个 Gin 引擎
	routes.HealthCheckInit(r) // 调用你的应用程序的路由设置函数

	// 使用应用程序的路由处理请求
	r.ServeHTTP(rr, req)

	// 检查响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

}
