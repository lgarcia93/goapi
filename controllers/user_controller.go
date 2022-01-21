package controller

import (
	"fitgoapi/jwt"
	"fitgoapi/middleware"
	"fitgoapi/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	email := c.GetHeader("email")
	password := c.GetHeader("password")
	fcmToken := c.GetHeader("fcmToken")

	userCredentialsAreValid := userService.Login(email, password)

	if userCredentialsAreValid {

		token := jwt.JWTService().GenerateToken(email, true)

		_ = userService.UpdateFcmToken(email, fcmToken)

		c.Header("Authorization", fmt.Sprintf("Bearer %s", token))

		c.Status(http.StatusOK)

		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
}

func loadInstructors(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	user := c.Value("User").(model.User)

	pageable, err := userService.FindAllInstructorsByCityCode(user.ID, c.Param("cityCode"), offset, size)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	} else {
		c.JSON(http.StatusOK, pageable)

		return
	}
}

func fetchMyProfile(c *gin.Context) {

	user := c.Value("User").(model.User)

	c.JSON(http.StatusOK, user)
}

func updateUserProfile(c *gin.Context) {

	var updatedUser model.User

	c.BindJSON(&updatedUser)

	user := c.Value("User").(model.User)

	err := userService.UpdateUser(user.ID, updatedUser)

	if err != nil {

		c.AbortWithStatus(http.StatusConflict)
	}

	return

}

func connectToUser(c *gin.Context) {

	contactId, err := strconv.Atoi(c.Param("contactId"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	contactID64 := int64(contactId)

	user := c.Value("User").(model.User)

	err = userService.MakeConnection(user.ID, contactID64)

	if err != nil {
		fmt.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, true)
	}

	return
}

func getConnections(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	user := c.Value("User").(model.User)

	pageable, err := userService.LoadUserConnections(user.ID, page, size)

	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, pageable)
	}

	return
}

// InitUserController init
func InitUserController(router *gin.Engine) {

	public := router.Group("/user")
	{

		public.POST("/login", login)

	}

	privateUser := router.Group("/user")
	{

		privateUser.Use(middleware.JWTValidator())

		privateUser.GET("/:cityCode/instructors", loadInstructors)

	}

	private := router.Group("/profile")
	{
		private.Use(middleware.JWTValidator())

		private.GET("/myself", fetchMyProfile)

		private.PATCH("/update", updateUserProfile)

		private.POST("/:contactId/connect", connectToUser)

		private.GET("/connections", getConnections)

	}

}
