package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
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

func UpdateOrderRoute(c *gin.Context) {
	boxID := c.Param("BoxID")

	var boxIDInt, err = strconv.Atoi(boxID)
	if err != nil {
		// Handle the error if the conversion fails
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BoxID"})
		return
	}

	// Get the stops array from the request body
	var stops []string
	if err := c.BindJSON(&stops); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the order with the given BoxID
	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").UpdateOne(
		context.TODO(),
		bson.M{"boxid": boxIDInt},
		bson.M{"$set": bson.M{"packageroute.stops": stops, "status": "In Route"}},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		}
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)
}
