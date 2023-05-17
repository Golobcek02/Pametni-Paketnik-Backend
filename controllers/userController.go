package controllers

import (
	"backend/schemas"
	"backend/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LogIn(c *gin.Context) {
	var loginUser schemas.User

	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res schemas.User
	err := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.Background(), bson.M{"username": loginUser.Username, "password": utils.Hash(loginUser.Password)}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusForbidden, "Denied")
		return
	}
	fmt.Printf(res.Username)
	c.IndentedJSON(http.StatusOK, res)
}

func Register(c *gin.Context) {
	var registerUser schemas.User

	if err := c.BindJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result schemas.User
	res := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.Background(), bson.M{"username": registerUser.Username}).Decode(&result)
	if res == mongo.ErrNoDocuments {
		registerUser.Password = utils.Hash(registerUser.Password)
		result, err := utils.CheckBase().Database("PametniPaketnik").Collection("users").InsertOne(context.TODO(), registerUser)
		fmt.Println(err)
		fmt.Println(result.InsertedID)
		c.IndentedJSON(http.StatusOK, bson.M{"res": "Proceede"})
		return
	}

	c.IndentedJSON(http.StatusBadRequest, bson.M{"res": "Denied"})
}
