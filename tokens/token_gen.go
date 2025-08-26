package tokens

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patil-prathamesh/e-commerce-golang/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
	jwt.RegisteredClaims
}

var SECRET_KEY string

func InitSecretKey() {
	SECRET_KEY = os.Getenv("SECRET_KEY")
}

var UserData *mongo.Collection = database.UserData(database.Client, "users")

func TokenGenerator(userID string, email string, firstName string, lastName string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := &SignedDetails{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedToken, err = token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err = refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		msg = "token is expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.M{"_id": userObjectID}
	update := bson.M{
		"$set": bson.M{
			"refresh_token": signedRefreshToken,
			"updated_at":    time.Now(),
		},
	}

	result, err := UserData.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Panic(err)
	}

	if result.ModifiedCount > 0 {
		fmt.Println("Existing document updated")
	}
}
