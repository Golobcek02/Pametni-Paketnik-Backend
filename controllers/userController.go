package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Register function handles the user registration process.
func Register(c *gin.Context) {
	var registerUser schemas.User

	if err := c.BindJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result schemas.User
	res := utils.CheckBase().Database("drvocepalci").Collection("users").FindOne(context.Background(), bson.M{"username": registerUser.Username}).Decode(&result)
	if res == mongo.ErrNoDocuments {
		registerUser.Password = utils.Hash(registerUser.Password)
		result, err := utils.CheckBase().Database("drvocepalci").Collection("users").InsertOne(context.TODO(), registerUser)
		fmt.Println(err)
		fmt.Println(result.InsertedID)
		c.IndentedJSON(http.StatusOK, "Poceede")
		return
	}

	c.IndentedJSON(http.StatusBadRequest, "Denied")
}
