package main

import (
	"backend/endpoints"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	Router := gin.Default()
	Router.Use(cors.Default())

	endpoints.Router(Router)

	Router.Run("localhost:5551")
}
