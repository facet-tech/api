package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	fmt.Println("mpika man mou!")
	return func(c *gin.Context) {
		fmt.Println("MESA!",c.GetHeader("Authorization"))
		//const BEARER_SCHEMA = "Bearer"
		//authHeader := c.GetHeader("Authorization")
		//tokenString := authHeader[len(BEARER_SCHEMA):]
		//token, err := service.JWTAuthService().ValidateToken(tokenString)
		//
		//if token.Valid {
		//	claims := token.Claims.(jwt.MapClaims)
		//	fmt.Println(claims)
		//} else {
		//	fmt.Println(err)
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//}

	}
}