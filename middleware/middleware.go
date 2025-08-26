package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/patil-prathamesh/e-commerce-golang/tokens"
)

func Authentication(c *gin.Context) {
    clientToken := c.Request.Header.Get("Authorization")
    
    if clientToken == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
        c.Abort()
        return
    }

    if len(clientToken) > 7 && clientToken[:7]=="Bearer " {
        clientToken = clientToken[7:]
    }

    claims, msg := tokens.ValidateToken(clientToken)
    if msg != "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
        c.Abort()
        return
    }

    c.Set("email", claims.Email)
    c.Set("first_name", claims.FirstName)
    c.Set("last_name", claims.LastName)
    
    c.Next()
}