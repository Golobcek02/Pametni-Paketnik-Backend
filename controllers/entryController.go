package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserEntries(c *gin.Context) {
	var allOpenings []schemas.Entry
	var res schemas.User
	var boxids []int

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		panic(err)
	}

	error := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&res)
	s := strings.Split(res.UserBoxes, "||")

	//fmt.Print(error)
	//fmt.Print(res.Username)

	for _, str := range s {
		intVal, _ := strconv.Atoi(str)
		boxids = append(boxids, intVal)
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
		/*for _, str := range s {
		x, _ := strconv.Atoi(str)
		if x == elem.BoxId {*/
		allOpenings = append(allOpenings, elem)
		/*}
		}*/
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
