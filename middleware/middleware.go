package middleware

import (
	"fmt"
	"net/http"
	"strings"

	j "fitgoapi/jwt"

	"fitgoapi/repository"

	"github.com/gin-gonic/gin"
)

// JWTValidator jwt validator
func JWTValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		bearerTokenParts := strings.Split(authHeader, " ")

		if len(bearerTokenParts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := j.JWTService().ValidateToken(tokenString)

		if err != nil {

			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		if token.Valid {
			claims := token.Claims.(*j.AuthCustomClaims)

			user, err := repository.UserRepository{}.FetchUserByEmail(claims.UserName)

			if err != nil {

				c.AbortWithStatus(http.StatusConflict)

				return
			}

			c.Set("User", user)

		} else {
			fmt.Printf("%s \n", err)

			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
