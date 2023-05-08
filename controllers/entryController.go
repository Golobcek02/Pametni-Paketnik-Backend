package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserEntries(c *gin.Context) {
	var allOpenings []schemas.Entry
	var boxids []int
	var tid = c.Param("id")

	cur, _ := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: tid}})
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		boxids = append(boxids, elem.BoxId)
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
	var newOpening schemas.Entry

	if err := c.BindJSON(&newOpening); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("entries").InsertOne(context.TODO(), newOpening)
	if err != nil {
		c.IndentedJSON(http.StatusConflict, err)
	}
	fmt.Println(err)
	fmt.Println(result.InsertedID)
	c.IndentedJSON(http.StatusOK, "Proceede")
}
