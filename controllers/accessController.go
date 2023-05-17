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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAccess(c *gin.Context) {
	var requestData struct {
		UserID   string
		AccessId string
		BoxId    int
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var res schemas.Box
	error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").FindOne(context.TODO(), bson.D{{Key: "boxid", Value: requestData.BoxId}}).Decode(&res)

	if error == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	if res.OwnerId == str {
		var r schemas.User
		error := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.TODO(), bson.D{{Key: "username", Value: requestData.AccessId}}).Decode(&r)
		if error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
			return
		}

		res.AccessIds = append(res.AccessIds, r.ID)
		_, e := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(
			context.Background(),
			bson.M{"boxid": requestData.BoxId},
			bson.M{"$set": res},
			options.Update().SetUpsert(true),
		)
		fmt.Println(e)
		c.IndentedJSON(http.StatusOK, "Added")
		return

	} else {
		c.IndentedJSON(http.StatusBadRequest, "Error")
		return
	}
}

func RevokeAccess(c *gin.Context) {
	var requestData struct {
		UserID   string
		AccessId string
		BoxId    int
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var res schemas.Box
	error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").FindOne(context.TODO(), bson.D{{Key: "boxid", Value: requestData.BoxId}}).Decode(&res)

	if error == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	if res.OwnerId == str {
		var r schemas.User
		error := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.TODO(), bson.D{{Key: "username", Value: requestData.AccessId}}).Decode(&r)
		if error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
			return
		}

		var ret []primitive.ObjectID
		for _, v := range res.AccessIds {
			if v != r.ID {
				ret = append(ret, v)
			}
		}

		res.AccessIds = ret
		_, e := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(
			context.Background(),
			bson.M{"boxid": requestData.BoxId},
			bson.M{"$set": res},
			options.Update().SetUpsert(true),
		)
		fmt.Println(e)
		c.IndentedJSON(http.StatusOK, "Removed!")
		return

	} else {
		c.IndentedJSON(http.StatusBadRequest, "Error")
		return
	}
}

func CheckAccess(c *gin.Context) {
	var requestData struct {
		UserID string
		BoxId  int
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var r schemas.User
	error := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.TODO(), bson.D{{Key: "username", Value: requestData.UserID}}).Decode(&r)
	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	var allBoxes []schemas.Box
	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{}})
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

	for _, v := range allBoxes {
		for _, z := range v.AccessIds {
			if z == r.ID {
				c.IndentedJSON(http.StatusOK, "Autorised!")
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, "Error")
	return
}
