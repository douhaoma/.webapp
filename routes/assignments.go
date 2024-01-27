package routes

import (
	".webapp/controllers"
	".webapp/middlewares"
	"github.com/gin-gonic/gin"
)

func AssignmentRoutesInit(r *gin.Engine) {
	assignmentGroup := r.Group("/v1/assignments", middlewares.InitBasicAuth)
	{
		assignmentGroup.Use(middlewares.CheckMethodNotAllowed(r))
		assignmentGroup.GET("", controllers.GetAssignmentsList)
		assignmentGroup.POST("", controllers.CreateAssignment)
		//assignmentGroup.GET("/", controllers.GetAssignmentById)       //用c.Query()!!
		assignmentGroup.GET("/:id", controllers.GetAssignmentById) //用c.Param()
		//assignmentGroup.DELETE("/", controllers.DeleteAssignmentById)
		assignmentGroup.DELETE("/:id", controllers.DeleteAssignmentById)
		//assignmentGroup.PUT("/", controllers.UpdateAssignment)
		assignmentGroup.PUT("/:id", controllers.UpdateAssignment)
		assignmentGroup.POST("/:id/submission", controllers.SubmitAssignment)
	}

}
