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
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
	}

	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist."})
	}

	count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
	token, refreshToken, _ := tokens.TokenGenerator(user.Email, user.FirstName, user.LastName)
	user.Token = token
	user.RefreshToken = refreshToken
	user.UserCart = make([]models.ProductUser, 0)
	user.AddressDetails = make([]models.Address, 0)
	user.OrderStatus = make([]models.Order, 0)

	_, insertErr := UserCollection.InsertOne(ctx, user)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user did not get created"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User creted successfully",
		"user_id": user.ID.Hex(),
	})

}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
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

	token, refreshToken, _ := tokens.TokenGenerator(foundUser.Email, foundUser.FirstName, foundUser.LastName)

	tokens.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())

	c.JSON(http.StatusFound, foundUser)
}

func ProductViewerAdmin(c *gin.Context) {}

func SearchProduct(c *gin.Context) {
	var productList []models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	cursor, err := ProductCollection.Find(ctx,bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"something went wrong, please try after some time"})
		return
	}

	err = cursor.All(ctx, &productList)
	if err != nil {
		log.Println(err)
		c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": productList,
		"count": len(productList),
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

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":query, "$options": "i"}})

	if err != nil {
		c.JSON(404, "something went wrong")
		return
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &searchedProducts)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error":"invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": searchedProducts,
		"count": len(searchedProducts),
	})
}