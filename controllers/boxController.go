package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AddUserBox(c *gin.Context) {
	var requestData struct {
		UserID     string `json:"user_id"`
		SmartBoxID string `json:"smartbox_id"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(requestData.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user schemas.User
	err = utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUserBoxes := user.UserBoxes
	if newUserBoxes != "" {
		newUserBoxes += "||"
	}
	newUserBoxes += requestData.SmartBoxID

	updateResult, err := utils.CheckBase().Database("PametniPaketnik").Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"userboxes": newUserBoxes}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or smartbox ID already in userboxes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Smartbox ID successfully appended to userboxes"})
}
