package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	Router := gin.Default()
	Router.Use(cors.Default())

	Router.Run("localhost:5551")
}
