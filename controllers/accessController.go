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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAccess(c *gin.Context) {
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
	fmt.Println(error)

	var finStr = res.AccessIds + requestData.AccessId

	_, e := utils.CheckBase().Database("PametniPaketnik").Collection("access").UpdateOne(context.Background(), bson.M{"ownerid": str}, bson.D{{Key: "$set", Value: bson.D{{Key: "accessids", Value: finStr}}}}, options.Update().SetUpsert(true))

	if e != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "AccessSuccesfully added!"})
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

	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	res.AccessIds = utils.RewokeAccess(res.AccessIds, requestData.AccessId)

	_, e := utils.CheckBase().Database("PametniPaketnik").Collection("access").UpdateOne(context.Background(), bson.M{"ownerid": str}, bson.D{{Key: "$set", Value: bson.D{{Key: "accessids", Value: res.AccessIds}}}})

	if e != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Clearance has been rewoked!")
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

	if utils.GetMatch(res.AccessIds, requestData.AccessId) {
		c.IndentedJSON(http.StatusOK, "Allowed")
	} else {
		c.IndentedJSON(http.StatusForbidden, "Denied")
	}

}
