package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/e-commerce-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress(c *gin.Context) {}

func EditAddress(c *gin.Context) {}

func EditWorkAddress(c *gin.Context) {}

func DeleteAddress(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	addresses := []models.Address{}
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	filter := bson.M{"_id": userObjectId}
	update := bson.M{"$set": bson.M{"address_details": addresses}}

	result, err := UserCollection.UpdateOne(ctx,filter,update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted successfully"})
}
