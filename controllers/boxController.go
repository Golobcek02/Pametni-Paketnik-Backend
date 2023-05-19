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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddUserBox(c *gin.Context) {
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

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		tmp, _ := strconv.Atoi(requestData.SmartBoxID)
		if elem.BoxId == tmp && elem.OwnerId.Hex() == "" {
			_, error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.Background(),
				bson.D{{Key: "boxid", Value: elem.BoxId}},
				bson.D{{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: str},
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
	//str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var box schemas.Box
	var temp []primitive.ObjectID
	box.BoxId, _ = strconv.Atoi(requestData.SmartBoxID)
	box.Latitude = requestData.Lat
	box.Longitude = requestData.Lon
	box.OwnerId = str
	box.AccessIds = temp
	fmt.Println(box)

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").InsertOne(context.Background(), box)
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

func ClearBoxOwner(c *gin.Context) {
	boxid := c.Param("id")
	boxIdInt, _ := strconv.Atoi(boxid)

	filter := bson.D{{Key: "boxid", Value: boxIdInt}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ownerid", Value: ""}}}}

	res, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.TODO(), filter, update)
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

// preimenuj v get user boxes and acesses, naredi dodatno da samo boxe vrne
func GetUserBoxesAndAccesses(c *gin.Context) {
	var allBoxes []schemas.Box
	var usrid = c.Param("id")
	str, _ := primitive.ObjectIDFromHex(usrid)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: str}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	var usernames [][]string
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		if len(elem.AccessIds) > 0 {
			var boxUsernames []string
			for _, id := range elem.AccessIds {
				user := schemas.User{}
				err := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
				if err == mongo.ErrNoDocuments {
					continue // user not found, skip to next id
				} else if err != nil {
					log.Fatal(err)
				}
				boxUsernames = append(boxUsernames, user.Username)
			}
			usernames = append(usernames, boxUsernames)
		} else {
			usernames = append(usernames, []string{})
		}

		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
	}
	obj := bson.M{"allBoxes": allBoxes, "usernames": usernames}
	c.IndentedJSON(http.StatusOK, obj)
}

func GetUserBoxes(c *gin.Context) {
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
			log.Fatal(err)
		}

		// Ignore the AccessIds for this request
		elem.AccessIds = nil

		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
	}
	obj := bson.M{"allBoxes": allBoxes}
	c.IndentedJSON(http.StatusOK, obj)
}
