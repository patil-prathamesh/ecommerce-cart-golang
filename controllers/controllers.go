package controllers

import "github.com/gin-gonic/gin"

func HashPassword(password string) string {}

func VerifyPassword(userPassword, givenPassword string) (bool,string) {}

func SignUp(c *gin.Context) {}

func Login(c *gin.Context) {}

func ProductViewerAdmin(c *gin.Context) {}

func SearchProduct(c *gin.Context) {}

func SearchProductByQuery(c *gin.Context) {}

func NewApplication() {}
