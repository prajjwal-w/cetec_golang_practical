package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prajjwal-w/cetec_golang_practical/controller"
)

func Routes(routes *gin.Engine) {
	routes.GET("/person/:person_id/info", controller.GetPerson())
	routes.POST("/person/create", controller.CreatePerson())
}
