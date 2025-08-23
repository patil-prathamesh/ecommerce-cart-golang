package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/e-commerce-golang/database"
	"github.com/patil-prathamesh/e-commerce-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ProdCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewApplication(ProdCollection, UserCollection *mongo.Collection) *Application {
	return &Application{ProdCollection, UserCollection}
}

func (app *Application) AddToCart(c *gin.Context) {
	productQueryId := c.Query("product_id")
	if productQueryId == "" {
		log.Println("product id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}

	userQueryId := c.Query("user_id")
	if userQueryId == "" {
		log.Println("user id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	productId, err := primitive.ObjectIDFromHex(productQueryId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = database.AddProductToCart(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "successfully added to the cart"})
}

func (app *Application) RemoveItem(c *gin.Context) {
	productQueryId := c.Query("product_id")
	if productQueryId == "" {
		log.Println("product id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}

	userQueryId := c.Query("user_id")
	if userQueryId == "" {
		log.Println("user id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	productId, err := primitive.ObjectIDFromHex(productQueryId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err = database.RemoveCartItem(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(200, "successfully removed item from cart")
}

func GetItemFromCart(c *gin.Context) {
	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
		return
	}

	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var filledCart models.User

	err := UserCollection.FindOne(ctx, bson.M{"_id": userObjectId}).Decode(&filledCart)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	filterMatch := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userObjectId}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user_cart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$user_cart.price"}}}}}}

	cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unwind, grouping})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var listing []bson.M
	if err := cursor.All(ctx, &listing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": listing[0]["total"]})
}

func (app *Application) BuyFromCart(c *gin.Context) {
	userQueryId := c.Query("user_id")
	if userQueryId == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := database.BuyItemFromCart(ctx, app.UserCollection, userQueryId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(200, "successfully placed the order")
}

func (app *Application) InstantBuy(c *gin.Context) {
	productQueryId := c.Query("product_id")
	if productQueryId == "" {
		log.Println("product id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}

	userQueryId := c.Query("user_id")
	if userQueryId == "" {
		log.Println("user id is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	productId, err := primitive.ObjectIDFromHex(productQueryId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = database.InstantBuyer(ctx, app.ProdCollection, app.UserCollection, productId, userQueryId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(200, "successfully placed the order")
}
