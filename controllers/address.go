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

func AddAddress(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	var address models.Address

	address.AddressID = primitive.NewObjectID()

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": userObjectId}
	update := bson.M{"$push": bson.M{"address_details": address}}

	result, err := UserCollection.UpdateOne(ctx, filter, update)

	if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
        "message": "address added successfully",
        "address_id": address.AddressID.Hex(),
    })
}

func EditHomeAddress(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	var editAddress models.Address

	if err := c.ShouldBindJSON(&editAddress); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

	filter := bson.M{"_id" : userObjectId}
	update := bson.M{"$set" : bson.M{
		"address_details.0.house": editAddress.House,
		"address_details.0.street": editAddress.Street,
		"address_details.0.city": editAddress.City,
		"address_details.0.pincode": editAddress.Pincode,
	}}

	result, err := UserCollection.UpdateOne(ctx, filter, update)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update home address"})
        return
    }

	if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "home address updated successfully"})
}

func EditWorkAddress(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	var editAddress models.Address

	if err := c.ShouldBindJSON(&editAddress); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

	filter := bson.M{"_id" : userObjectId}
	update := bson.M{"$set" : bson.M{
		"address_details.1.house": editAddress.House,
		"address_details.1.street": editAddress.Street,
		"address_details.1.city": editAddress.City,
		"address_details.1.pincode": editAddress.Pincode,
	}}

	result, err := UserCollection.UpdateOne(ctx, filter, update)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update work address"})
        return
    }

	if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "work address updated successfully"})
}

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
