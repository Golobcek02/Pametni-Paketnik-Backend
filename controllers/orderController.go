package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertOrder(c *gin.Context) {
	var order schemas.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//order.ID = primitive.NewObjectID()

	// Insert the order into the database
	insertResult, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order"})
		return
	}

	// Return the inserted order ID
	c.JSON(http.StatusOK, gin.H{"orderId": insertResult.InsertedID})
}
