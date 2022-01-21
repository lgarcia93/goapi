package controller

import (
	s "fitgoapi/service"

	"github.com/gin-gonic/gin"
)

var configService s.IConfigService = s.ConfigService{}

func loadConfigs(c *gin.Context) {
	configs, err := configService.LoadConfigs()
	if err != nil {
		c.Status(400)
	} else {
		c.JSON(200, configs)
	}
}

// InitUserController init
func InitConfigController(router *gin.Engine) {

	public := router.Group("/config")
	{
		public.GET("/setup", loadConfigs)

	}
}
