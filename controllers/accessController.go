package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAccess(c *gin.Context) {
	var requestData struct {
		UserID     string
		SmartBoxID string
		Lat        float64
		Lon        float64
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(requestData)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.Background(), bson.D{{}})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		tmp, _ := strconv.Atoi(requestData.SmartBoxID)
		if elem.BoxId == tmp && elem.OwnerId == "" {
			_, error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.Background(),
				bson.D{{Key: "boxid", Value: elem.BoxId}},
				bson.D{{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: requestData.UserID},
					{Key: "latitude", Value: requestData.Lat},
					{Key: "longitude", Value: requestData.Lon},
				}}},
				options.Update().SetUpsert(true))

			if error != nil {
				panic(error)
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "Box successfully updated!"})
			return
		}
	}

	var box schemas.Box
	box.BoxId, _ = strconv.Atoi(requestData.SmartBoxID)
	box.Latitude = requestData.Lat
	box.Longitude = requestData.Lon
	box.OwnerId = requestData.UserID
	fmt.Println(box)

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").InsertOne(context.Background(), box)
	fmt.Println(err)
	fmt.Println(result.InsertedID)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Box successfully inserted!"})
}

func RewokeAccess(c *gin.Context) {
	var requestData struct {
		UserID   string
		AccessId string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var res schemas.Access
	error := utils.CheckBase().Database("PametniPaketnik").Collection("access").FindOne(context.Background(), bson.M{"ownerid": str}).Decode(&res)

	res, err := utils.CheckBase().Database("PametniPaketnik").Collection("access").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error while clearing the owner of this box"})
		return
	}

	if res.ModifiedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Box not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Owner of the box successfully cleared"})
}

func CheckAccess(c *gin.Context) {
	var requestData struct {
		UserID   string
		AccessId string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var res schemas.Access
	error := utils.CheckBase().Database("PametniPaketnik").Collection("access").FindOne(context.Background(), bson.M{"ownerid": str}).Decode(&res)

	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	if utils.GetMatch(requestData.AccessId) {
		c.IndentedJSON(http.StatusOK, "Allowed")
	} else {
		c.IndentedJSON(http.StatusForbidden, "Denied")
	}

}
