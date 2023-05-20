package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertPackageRoutes(c *gin.Context) {
	var requestData struct {
		route string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	BoxIds := strings.Split(requestData.route, "|")

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, "No boxes found")
	}

	var box []schemas.Box
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range BoxIds {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if i == elem.BoxId {
				box = append(box, elem)
			}
		}
	}

	curr, e := utils.CheckBase().Database("PametniPaketnik").Collection("orders").Find(context.TODO(), bson.D{{}})
	if e == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, "No orders found")
	}

	var ord []schemas.Order
	for curr.Next(context.TODO()) {
		var elem schemas.Order
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range BoxIds {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if i == elem.BoxID {
				ord = append(ord, elem)
			}
		}

	}

	centralStation := strings.Split(requestData.route, ":")
	var station schemas.Order
	station.PageUrl = centralStation[0]
	station.BoxID = 000

	var packageRoute schemas.PackageRoutes
	packageRoute.Orders = append(packageRoute.Orders, station)
	packageRoute.Stops = append(packageRoute.Stops, centralStation[1])
	for _, v := range box {
		for _, z := range ord {
			if v.BoxId == z.BoxID {
				lat := strconv.FormatFloat(v.Latitude, 'f', 2, 64)
				lon := strconv.FormatFloat(v.Longitude, 'f', 2, 64)
				packageRoute.Orders = append(packageRoute.Orders, z)
				packageRoute.Stops = append(packageRoute.Stops, lat+", "+lon)
			}
		}
	}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("packageRoutes").InsertOne(context.TODO(), packageRoute)
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

func PopFirstStop(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	filter := bson.M{"boxid": id}
	update := bson.M{"$pop": bson.M{"packageroute.stops": -1}} // Use -1 to pop the first element

	_, err = utils.CheckBase().Database("PametniPaketnik").Collection("orders").UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
