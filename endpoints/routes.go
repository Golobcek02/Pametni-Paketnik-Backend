package endpoints

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Router(Router *gin.Engine) {
	Router.POST("/register", controllers.Register)
	Router.POST("/login", controllers.LogIn)
	Router.POST("/newEntry", controllers.InsertNewEntry)
	Router.GET("/getEntries/:id", controllers.GetUserEntries)
	Router.POST("/addUserBox", controllers.AddUserBox)
	Router.DELETE("/removeEntry", controllers.RemoveEntry)
	Router.DELETE("/removeBox", controllers.RemoveBox)
}
