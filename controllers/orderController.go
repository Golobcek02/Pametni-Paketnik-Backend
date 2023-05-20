package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertOrder(c *gin.Context) {
	var order schemas.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the order into the database
	insertResult, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order"})
		return
	}

	// Return the inserted order ID
	c.JSON(http.StatusOK, gin.H{"orderId": insertResult.InsertedID})
}

func GetUserOrders(c *gin.Context) {
	var allBoxes []schemas.Box
	var usrid = c.Param("id")
	str, _ := primitive.ObjectIDFromHex(usrid)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: str}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Error")
			return
		}

		elem.AccessIds = nil

		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
		return
	}

	boxIDs := make([]int, len(allBoxes))
	for i, box := range allBoxes {
		boxIDs[i] = box.BoxId
	}

	ordersCur, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").Find(context.TODO(), bson.M{
		"boxid": bson.M{"$in": boxIDs},
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
		return
	}

	var orders []schemas.Order
	if err := ordersCur.All(context.TODO(), &orders); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
		return
	}

	obj := bson.M{"orders": orders}
	c.IndentedJSON(http.StatusOK, obj)
}

func UpdateOrderStatus(c *gin.Context) {

	var requestData struct {
		BoxId  string
		Status string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order []schemas.Order
	curr, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").Find(context.TODO(), bson.M{"boxid": requestData.BoxId})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for curr.Next(context.TODO()) {
		var elem schemas.Order
		err := curr.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		order = append(order, elem)
	}

	for _, v := range order {
		_, error := utils.CheckBase().Database("PametniPaketnik").Collection("orders").UpdateOne(context.TODO(), bson.D{{"_id", v.ID}}, bson.D{{"$set", bson.D{{"Status", requestData.Status}}}}, options.Update().SetUpsert(false))
		if error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"code": "Success"})
}
