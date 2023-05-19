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
	Router.DELETE("/removeEntry/:id", controllers.RemoveEntry)
	Router.DELETE("/removeBox/:id", controllers.RemoveBox)
	Router.GET("/getUserBoxesAndAccesses/:id", controllers.GetUserBoxesAndAccesses)
	Router.GET("/getUserBoxes/:id", controllers.GetUserBoxes)
	Router.PUT("/clearBox/:id", controllers.ClearBoxOwner)
	Router.POST("/addAccessToUser", controllers.AddAccess)
	Router.POST("/revokeAccessToUser", controllers.RevokeAccess)
	Router.POST("/checkAccessOfUser", controllers.CheckAccess)
	Router.POST("/addPackageRoute", controllers.InsertPackageRoutes)
	Router.POST("/addOrder", controllers.InsertOrder)
	Router.GET("/getUserOrders/:id", controllers.GetUserOrders)
	Router.POST("/updateOrderRoute/:BoxID", controllers.UpdateOrderRoute)
	Router.POST("/popFirstStop/:id", controllers.PopFirstStop)
	Router.POST("/login2fa", controllers.Temp)
}
