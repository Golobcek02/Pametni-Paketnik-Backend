package endpoints

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Router(Router *gin.Engine) {
	Router.POST("/register", controllers.Register)
	Router.POST("/login", controllers.LogIn)
	Router.POST("/newEntry/:id", controllers.InsertNewEntry)
	Router.GET("/getEntries", controllers.GetUserEntries)
}
