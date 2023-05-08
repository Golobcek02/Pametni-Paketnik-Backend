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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddUserBox(c *gin.Context) {
	var requestData struct {
		UserID     string
		SmartBoxID string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.Background(), bson.D{{Key: "ownerid", Value: requestData.UserID}})
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
		if elem.BoxId == tmp {
			_, error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.Background(), bson.D{{Key: "boxid", Value: elem.BoxId}}, bson.D{{Key: "$set", Value: bson.D{{Key: "ownerid", Value: "0"}}}}, options.Update().SetUpsert(true))
			if error != nil {
				panic(error)
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "Box successfully updated!"})
			return
		}
	}
	var box schemas.Box
	box.BoxId, _ = strconv.Atoi(requestData.SmartBoxID)
	box.OwnerId = requestData.UserID
	fmt.Println(box)
	fmt.Println(requestData)

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").InsertOne(context.TODO(), box)
	fmt.Println(err)
	fmt.Println(result.InsertedID)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Box successfully inserted!"})
}

func RemoveBox(c *gin.Context) {
	var boxid = c.Param("id")
	boxIdInt, _ := strconv.Atoi(boxid)

	res, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").DeleteOne(context.TODO(), bson.D{{Key: "boxid", Value: boxIdInt}})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Error while deleting this box")
	}

	fmt.Print(res.DeletedCount)
	c.IndentedJSON(http.StatusOK, "successfully deleted")
}

func GetUserBoxes(c *gin.Context) {
	var allBoxes []schemas.Box
	var usrid = c.Param("id")

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: usrid}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
	}
	c.IndentedJSON(http.StatusOK, allBoxes)
}
