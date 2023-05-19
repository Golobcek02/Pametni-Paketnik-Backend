package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserEntries(c *gin.Context) {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var boxes []schemas.Box
	boxIds := []int{}
	boxFilter := bson.M{"loggerid": objectId}
	boxCursor, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), boxFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find boxes"})
		return
	}
	for boxCursor.Next(context.Background()) {
		var box schemas.Box
		if err := boxCursor.Decode(&box); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode box"})
			return
		}
		boxes = append(boxes, box)
		boxIds = append(boxIds, box.BoxId)
	}
	if err := boxCursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over boxes"})
		return
	}

	var entries []schemas.Entry
	entryFilter := bson.M{"boxid": bson.M{"$in": boxIds}}
	entryCursor, err := utils.CheckBase().Database("PametniPaketnik").Collection("entries").Find(context.TODO(), entryFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find entries"})
		return
	}
	for entryCursor.Next(context.Background()) {
		var entry schemas.Entry
		if err := entryCursor.Decode(&entry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode entry"})
			return
		}
		entries = append(entries, entry)
	}
	if err := entryCursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func InsertNewEntry(c *gin.Context) {
	var newEntry schemas.Entry

	if err := c.BindJSON(&newEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("entries").InsertOne(context.TODO(), newEntry)
	if err != nil {
		c.IndentedJSON(http.StatusConflict, err)
		return
	}
	fmt.Println(err)
	fmt.Println(result.InsertedID)
	c.IndentedJSON(http.StatusOK, "Proceede")
}

func RemoveEntry(c *gin.Context) {
	entryId := c.Param("id")
	str, err := primitive.ObjectIDFromHex(entryId)
	filter := bson.M{"_id": str}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("entries").DeleteOne(context.TODO(), filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Entry not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}
