package routes

import (
	".webapp/controllers"
	"github.com/gin-gonic/gin"
)

func HealthCheckInit(r *gin.Engine) {
	r.GET("/healthz", controllers.HealthCheck)
}
