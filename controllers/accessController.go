package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
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

	if res.OwnerId == requestData.UserID {
		var r schemas.User
		error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").FindOne(context.TODO(), bson.D{{Key: "boxid", Value: requestData.BoxId}}).Decode(&r)
		if error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
			return
		}

		_, e := utils.CheckBase().Database("PametniPaketnik").Collection("access").UpdateOne(
			context.Background(),
			bson.M{"ownerid": str},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: str},
					{Key: "accessids", Value: bson.D{{Key: "$concat", Value: bson.D{{Key: "$accessids", Value: " " + r.ID.String()}}}}},
					{Key: "boxid", Value: requestData.BoxId},
				}},
			},
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

func RewokeAccess(c *gin.Context) {
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

	if res.OwnerId == requestData.UserID {
		var r schemas.User
		error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").FindOne(context.TODO(), bson.D{{Key: "boxid", Value: requestData.BoxId}}).Decode(&r)
		if error != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
			return
		}

		_, e := utils.CheckBase().Database("PametniPaketnik").Collection("access").UpdateOne(
			context.Background(),
			bson.M{"ownerid": str},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: str},
					{Key: "accessids", Value: bson.D{{Key: "$regexReplace", Value: bson.M{"input": "$accessids", "find": requestData.AccessId, "replacement": "", "options": "i"}}}},
					{Key: "boxid", Value: requestData.BoxId},
				}},
			},
			options.Update().SetUpsert(true),
		)
		fmt.Println(e)
		c.IndentedJSON(http.StatusOK, "Revoked!")
		return

	} else {
		c.IndentedJSON(http.StatusBadRequest, "Error")
		return
	}
}

func CheckAccess(c *gin.Context) {
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

	var res schemas.Access
	error := utils.CheckBase().Database("PametniPaketnik").Collection("access").FindOne(context.Background(), bson.M{"ownerid": str}).Decode(&res)

	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	if utils.GetMatch(res.AccessIds, requestData.AccessId) {
		c.IndentedJSON(http.StatusOK, "Allowed")
	} else {
		c.IndentedJSON(http.StatusForbidden, "Denied")
	}

}

func GetAllAccess(c *gin.Context) {
	var requestData struct {
		UserID   string
		AccessId string
		BoxId    string
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

	if utils.GetMatch(res.AccessIds, requestData.AccessId) {
		c.IndentedJSON(http.StatusOK, "Allowed")
	} else {
		c.IndentedJSON(http.StatusForbidden, "Denied")
	}

}
