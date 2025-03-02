package controllers

import (
	"room-service/clients"
	"room-service/constants"
	"room-service/controllers"
	"room-service/middlewares"

	"github.com/gin-gonic/gin"
)

type TimeRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type ITimeRoute interface {
	Run()
}

func NewTimeRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) ITimeRoute {
	return &TimeRoute{controller: controller, group: group, client: client}
}

func (t *TimeRoute) Run() {
	group := t.group.Group("/time")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, t.client),
		t.controller.GetTime().GetAll)

	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, t.client),
		t.controller.GetTime().GetByUUID)

	group.POST("", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, t.client),
		t.controller.GetTime().Create)
}
