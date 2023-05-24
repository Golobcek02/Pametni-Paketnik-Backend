package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertPackageRoutes(c *gin.Context) {
	var requestData struct {
		Route string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	BoxIds := strings.Split(requestData.Route, "|")

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
		box = append(box, elem)
	}

	curr, e := utils.CheckBase().Database("PametniPaketnik").Collection("orders").Find(context.TODO(), bson.D{{}})
	if e == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, "No orders found")
	}

	var ord []schemas.Order
	for curr.Next(context.TODO()) {
		var element schemas.Order
		err := curr.Decode(&element)
		if err != nil {
			log.Fatal(err)
		}
		ord = append(ord, element)
	}

	centralStation := strings.Split(BoxIds[0], ":")
	var packageRoute schemas.PackageRoutes
	var zeroObjectID primitive.ObjectID
	packageRoute.Orders = append(packageRoute.Orders, zeroObjectID)
	packageRoute.Stops = append(packageRoute.Stops, centralStation[1])
	var idarr []primitive.ObjectID

	for _, v := range box {
		for _, z := range ord {
			for _, g := range BoxIds {
				i, _ := strconv.Atoi(g)
				if v.BoxId == z.BoxID && z.Status == "Pending" && i == z.BoxID {
					fmt.Println("neke")
					lat := strconv.FormatFloat(v.Latitude, 'f', 8, 64)
					lon := strconv.FormatFloat(v.Longitude, 'f', 8, 64)
					packageRoute.Orders = append(packageRoute.Orders, z.ID)
					packageRoute.Stops = append(packageRoute.Stops, lat+", "+lon)
					idarr = append(idarr, z.ID)
				}
			}
		}
	}

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("packageRoutes").InsertOne(context.TODO(), packageRoute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert package routes"})
		return
	}

	rs, err := utils.CheckBase().Database("PametniPaketnik").Collection("orders").UpdateMany(context.Background(), bson.M{"_id": bson.M{"$in": idarr}}, bson.M{"$set": bson.M{"status": "In Route"}})
	if err != nil {
		panic(err)
	}

	fmt.Println(rs)
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

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var packageRoute schemas.PackageRoutes
	err = utils.CheckBase().Database("PametniPaketnik").Collection("packageRoutes").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&packageRoute)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PackageRoute not found"})
		return
	}

	var zeroObjectID primitive.ObjectID

	i := 1

	if packageRoute.Orders[0] == zeroObjectID {
		packageRoute.Orders = packageRoute.Orders[i:]
		packageRoute.Stops = packageRoute.Stops[i:]
	} else {
		for i < len(packageRoute.Stops) && packageRoute.Stops[0] == packageRoute.Stops[i] {
			i++
		}

		_, err = utils.CheckBase().Database("PametniPaketnik").Collection("orders").UpdateMany(
			context.TODO(),
			bson.M{"_id": packageRoute.Orders[0]},
			bson.M{"$set": bson.M{"status": "Completed"}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PackageRoute"})
			return
		}

		var ord schemas.Order
		err = utils.CheckBase().Database("PametniPaketnik").Collection("orders").FindOne(context.TODO(), bson.M{"_id": packageRoute.Orders[0]}).Decode(&ord)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PackageRoute"})
			return
		}

		var entry schemas.Entry
		entry.BoxId = ord.BoxID
		entry.DeliveryId = 1
		entry.EntryType = "orderCompleted"
		entry.Latitude = 0
		entry.Longitude = 0
		entry.LoggerId = zeroObjectID
		entry.TimeAccessed = time.Now().Unix()

		_, err = utils.CheckBase().Database("PametniPaketnik").Collection("entries").InsertOne(context.Background(), entry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PackageRoute"})
			return
		}

		packageRoute.Orders = packageRoute.Orders[i:]
		packageRoute.Stops = packageRoute.Stops[i:]

		rs, _ := utils.CheckBase().Database("PametniPaketnik").Collection("orders").Find(context.Background(), bson.M{"orders": bson.M{"$in": packageRoute.Orders}})

		var entries []schemas.Entry
		for rs.Next(context.TODO()) {
			var element schemas.Order
			err := rs.Decode(&element)
			if err != nil {
				log.Fatal(err)
			}
			entry.BoxId = element.BoxID
			entry.DeliveryId = 2
			entry.EntryType = "oneStopCloser"
			entry.Latitude = 0
			entry.Longitude = 0
			entry.LoggerId = zeroObjectID
			entry.TimeAccessed = time.Now().Unix()
			entries = append(entries, entry)
		}

		var docs []interface{}
		for _, entry := range entries {
			docs = append(docs, entry)
		}
		_, err := utils.CheckBase().Database("PametniPaketnik").Collection("entries").InsertMany(context.Background(), docs)
		if err != nil {
			panic(err)
		}
	}

	_, err = utils.CheckBase().Database("PametniPaketnik").Collection("packageRoutes").UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"orders": packageRoute.Orders, "stops": packageRoute.Stops}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PackageRoute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
