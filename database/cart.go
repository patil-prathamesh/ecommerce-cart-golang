package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/patil-prathamesh/e-commerce-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productID})

	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchFromDb.All(ctx, &productCart)

	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.M{"_id": userObjectID}
	update := bson.M{"$push": bson.M{"user_cart": bson.M{"$each": productCart}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.M{"_id": userObjectID}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": productID}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var user models.User

	err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)

	if err != nil {
		log.Println(err)
		return ErrCantGetItem
	}
	if len(user.UserCart) == 0 {
		return errors.New("cart is empty")
	}

	var total uint64

	for _, v := range user.UserCart {
		total += v.Price
	}

	orderID := primitive.NewObjectID()

	newOrder := models.Order{
		OrderID:       orderID,
		OrderCart:     user.UserCart,
		OrderedAt:     time.Now(),
		Price:         total,
		Discount:      0,
		PaymentMethod: models.Payment{Digital: false, COD: true},
	}

	filter := bson.M{"_id": userObjectID}
	update := bson.M{
		"$push": bson.M{"orders": newOrder},
		"$set":  bson.M{"user_cart": []models.ProductUser{}}, // Clear cart after purchase
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var productDetails models.ProductUser
	err = prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&productDetails)

	if err != nil {
		log.Println(err)
	}

	orderDetails := models.Order{
		OrderID:       primitive.NewObjectID(),
		OrderCart:     []models.ProductUser{productDetails},
		OrderedAt:     time.Now(),
		Price:         productDetails.Price,
		Discount:      0,
		PaymentMethod: models.Payment{Digital: false, COD: true},
	}

	filter := bson.M{"_id": userObjectID}
	update := bson.M{"$push": bson.M{"orders": orderDetails}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter = bson.M{"_id": userObjectID}
	update = bson.M{"$push": bson.M{}}
	return nil
}
