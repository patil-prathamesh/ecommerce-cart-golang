package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/e-commerce-golang/controllers"
	"github.com/patil-prathamesh/e-commerce-golang/database"
	"github.com/patil-prathamesh/e-commerce-golang/middleware"
	"github.com/patil-prathamesh/e-commerce-golang/routes"
	"github.com/patil-prathamesh/e-commerce-golang/tokens"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	tokens.InitSecretKey()
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "products"), database.UserData(database.Client, "users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication)

	router.PUT("/addtocart", app.AddToCart)
	router.PUT("/removeitem", app.RemoveItem)
	router.POST("/cartcheckout", app.BuyFromCart)
	router.POST("/instantbuy", app.InstantBuy)

	log.Fatal(router.Run(":" + port))
}
