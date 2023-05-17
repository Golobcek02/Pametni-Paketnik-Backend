package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserEntries(c *gin.Context) {
	var allOpenings []schemas.Entry
	var boxids []int
	var tid = c.Param("id")

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: tid}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, allOpenings)
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		boxids = append(boxids, elem.BoxId)
	}

	if len(boxids) == 0 {
		c.IndentedJSON(http.StatusOK, allOpenings)
	}

	cur, error := utils.CheckBase().Database("PametniPaketnik").Collection("entries").Find(context.TODO(), bson.D{{Key: "boxid", Value: bson.D{{Key: "$in", Value: boxids}}}})
	if error != nil {
		panic(error)
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Entry
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		allOpenings = append(allOpenings, elem)
	}

	c.IndentedJSON(http.StatusOK, allOpenings)
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
