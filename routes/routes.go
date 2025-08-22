package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/e-commerce-golang/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp)
	incomingRoutes.POST("/users/login", controllers.Login)
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin)
	incomingRoutes.GET("/users/productView", controllers.SearchProduct)
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery)
}