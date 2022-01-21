package controller

import (
	"errors"
	"fitgoapi/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerUser(c *gin.Context) {

	var user model.User

	c.BindJSON(&user)

	createdUser, err := userService.CreateUser(user)

	if err != nil {
		fmt.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusConflict)

		return
	}

	c.JSON(http.StatusOK, createdUser)
}

func validateEmail(c *gin.Context) {

	email := c.DefaultQuery("email", "")

	isEmailValid, err := userService.ValidateEmail(email)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("erro ao validar dados"))
	} else if isEmailValid {
		c.JSON(http.StatusOK, true)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

	return
}

func InitAuthController(router *gin.Engine) {

	public := router.Group("/auth")
	{
		public.POST("/register", registerUser)
		public.GET("/validate", validateEmail)
	}
}
