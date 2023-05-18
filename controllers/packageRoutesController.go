package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertPackageRoutes(c *gin.Context) {
	var packageRoutes schemas.PackageRoutes

	if err := c.BindJSON(&packageRoutes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("packageRoutes").InsertOne(context.TODO(), packageRoutes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert package routes"})
		return
	}

	c.JSON(http.StatusOK, result.InsertedID)
}
