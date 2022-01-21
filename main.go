package main

import (
	controller "fitgoapi/controllers"
	database "fitgoapi/database"

	"os"

	gin "github.com/gin-gonic/gin"
)

func getPort() string {
	if os.Getenv("PORT") == "" {
		return "8080"
	}

	return os.Getenv("PORT")
}

func main() {

	database.ConfigDatabase()

	defer database.Connection.Close()

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	controller.InitUserController(r)
	controller.InitConfigController(r)
	controller.InitAuthController(r)
	controller.InitScheduleController(r)

	r.Run(":" + getPort())
}
