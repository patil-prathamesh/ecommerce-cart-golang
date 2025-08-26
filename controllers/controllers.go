package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/patil-prathamesh/e-commerce-golang/database"
	"github.com/patil-prathamesh/e-commerce-golang/models"
	"github.com/patil-prathamesh/e-commerce-golang/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or password is incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist."})
		return
	}

	count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number already in use"})
		return
	}

	password := HashPassword(user.Password)
	user.Password = password
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = primitive.NewObjectID()
	token, refreshToken, _ := tokens.TokenGenerator(user.ID.Hex(), user.Email, user.FirstName, user.LastName)
	user.RefreshToken = refreshToken
	user.UserCart = []models.ProductUser{}
	user.AddressDetails = []models.Address{}
	user.Order = []models.Order{}

	_, insertErr := UserCollection.InsertOne(ctx, user)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "User creted successfully",
		"user_id":      user.ID.Hex(),
		"access_token": token,
	})
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var foundUser models.User

	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
		return
	}

	IsPasswordValid, msg := VerifyPassword(user.Password, foundUser.Password)
	if !IsPasswordValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	token, refreshToken, _ := tokens.TokenGenerator(foundUser.ID.Hex(), foundUser.Email, foundUser.FirstName, foundUser.LastName)

	tokens.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())

	c.JSON(http.StatusFound, gin.H{
		"message": "user logged in",
		"access_token": token,
	})
}

func ProductViewerAdmin(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err !=  nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	count , _ := ProductCollection.CountDocuments(ctx, bson.M{"_id": product.ProductID})

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "product already exists"})
		return
	}

	inserted, err := ProductCollection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product inserted successfullly", "product_id": inserted.InsertedID})
}

func SearchProduct(c *gin.Context) {
	var productList []models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := ProductCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, please try after some time"})
		return
	}

	err = cursor.All(ctx, &productList)
	if err != nil {
		log.Println(err)
		c.AbortWithError(500, err) //
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": productList,
		"count":    len(productList),
	})
}

func SearchProductByQuery(c *gin.Context) {
	var searchedProducts []models.Product
	query := c.Query("name")

	if query == "" {
		log.Println("query is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search index"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": query, "$options": "i"}})

	if err != nil {
		c.JSON(404, "something went wrong")
		return
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &searchedProducts)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": searchedProducts,
		"count":    len(searchedProducts),
	})
}