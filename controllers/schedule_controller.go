package controller

import (
	"fitgoapi/middleware"
	"fitgoapi/model"
	s "fitgoapi/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var scheduleService s.IScheduleService = s.ScheduleService{}

func loadPendingRequests(c *gin.Context) {
	user := c.Value("User").(model.User)

	schedules, err := scheduleService.FindAllUserPendingSchedule(user.ID)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, schedules)
	}

	return
}

func loadUserSchedule(c *gin.Context) {

	week, _ := strconv.Atoi(c.DefaultQuery("week", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	dayOfWeek, _ := strconv.Atoi(c.DefaultQuery("dayOfWeek", "10"))
	user := c.Value("User").(model.User)

	if week == 0 {
		_, week = time.Now().ISOWeek()
	}

	if dayOfWeek == 0 {
		pageable, err := scheduleService.LoadUserSchedule(user.ID, week, page, size)

		if err != nil {
			fmt.Printf("%s", err.Error())
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, pageable)
	} else {
		pageable, err := scheduleService.LoadUserScheduleByDay(user.ID, week, dayOfWeek, page, size)

		if err != nil {
			fmt.Printf("%s", err.Error())
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, pageable)
	}

	return
}

func acceptOrRejectRequest(c *gin.Context) {
	accepted, err := strconv.ParseBool(c.DefaultQuery("accept", "true"))
	scheduleIDInt, err := strconv.Atoi(c.Param("scheduleId"))
	user := c.Value("User").(model.User)

	if err != nil {
		fmt.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	scheduleID := int64(scheduleIDInt)

	affected, err := scheduleService.AcceptOrRejectSchedule(user.ID, scheduleID, accepted)

	if err != nil {
		fmt.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, affected > 0)
	}

	return
}

func loadInstructorTimeline(c *gin.Context) {

	instructorID, err := strconv.Atoi(c.Param("userId"))

	instructorID64 := int64(instructorID)

	scheduleItems, err := scheduleService.LoadFullInstructorSchedule(instructorID64)

	if err != nil {
		fmt.Printf("%s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, scheduleItems)
	}

	return
}

func scheduleClass(c *gin.Context) {
	var items []model.ScheduleItemDTO

	c.BindJSON(&items)
	user := c.Value("User").(model.User)
	instructorID, err := strconv.Atoi(c.Param("userId"))
	instructorID64 := int64(instructorID)

	schedule, err := scheduleService.ScheduleClass(user.ID, instructorID64, items)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, schedule)
	}
}

func cancelScheduleEvent(c *gin.Context) {
	week, _ := strconv.Atoi(c.DefaultQuery("week", "0"))
	if week == 0 {
		_, week = time.Now().ISOWeek()
	}

	user := c.Value("User").(model.User)
	itemID, err := strconv.Atoi(c.Param("itemId"))
	itemID64 := int64(itemID)

	_, err = scheduleService.CancelScheduleEvent(user.ID, itemID64, week)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusAccepted, "")
	}
}

func changeScheduleEvent(c *gin.Context) {
	var incident model.IncidentDTO

	c.BindJSON(&incident)
	user := c.Value("User").(model.User)
	itemID, err := strconv.Atoi(c.Param("itemId"))
	itemID64 := int64(itemID)

	_, err = scheduleService.ChangeScheduleEvent(user.ID, itemID64, incident)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusAccepted, "")
	}
}

func acceptOrRejectChange(c *gin.Context) {
	var incident model.IncidentDTO

	c.BindJSON(&incident)

	changeID, err := strconv.Atoi(c.Param("idChange"))
	changeID64 := int64(changeID)

	/*itemID, err := strconv.Atoi(c.Param("itemId"))
	itemID64 := int64(itemID)*/

	answer, _ := strconv.ParseBool(c.DefaultQuery("response", "false"))

	user := c.Value("User").(model.User)

	err = scheduleService.AcceptOrRejectChange(user.ID, changeID64, answer)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusAccepted, changeID64)
	}
	return
}

func loadPendingChanges(c *gin.Context) {

	user := c.Value("User").(model.User)

	incidents, err := scheduleService.LoadPendingIncidents(user.ID, model.Cancellation)

	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
	} else {
		if incidents == nil {
			c.JSON(http.StatusNoContent, "") //TODO: gambeta de obj vazio
		} else {
			c.JSON(http.StatusOK, incidents)
		}
	}
}

// InitScheduleController init
func InitScheduleController(router *gin.Engine) {
	schedule := router.Group("schedule")
	{
		schedule.Use(middleware.JWTValidator())

		//done in app
		schedule.GET("/", loadUserSchedule)

		//done in app
		schedule.GET("/pending", loadPendingRequests)

		schedule.PUT("/pending/:scheduleId", acceptOrRejectRequest)

		//done in app
		schedule.GET("/instructor/:userId/unavailable", loadInstructorTimeline)

		//done in app
		schedule.POST("/instructor/:userId", scheduleClass)

		schedule.PUT("/item/:itemId/cancel", cancelScheduleEvent)
		schedule.PUT("/item/:itemId/change", changeScheduleEvent)
	}

	changes := router.Group("changes")
	{
		changes.Use(middleware.JWTValidator())
		changes.GET("/", loadPendingChanges)
		changes.PUT("/:idChange", acceptOrRejectChange)
	}
}
