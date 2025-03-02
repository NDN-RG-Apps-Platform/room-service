package routes

import (
	"room-service/clients"
	"room-service/constants"
	"room-service/controllers"
	"room-service/middlewares"

	"github.com/gin-gonic/gin"
)

type RoomRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRoomRoute interface {
	Run()
}

func NewRoomRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRoomRoute {
	return &RoomRoute{controller: controller, group: group, client: client}
}

func (r *RoomRoute) Run() {
	group := r.group.Group("/room")
	group.GET("", middlewares.AuthenticateWithoutToken(), r.controller.GetRoom().GetAllWithoutPagination)
	group.GET("/:uuid", middlewares.AuthenticateWithoutToken(), r.controller.GetRoom().GetByUUID)
	group.Use(middlewares.Authenticate())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
		constants.Lecture,
		constants.Student,
	}, r.client),
		r.controller.GetRoom().GetAllWithPagination)

	group.POST("", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoom().Create)

	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoom().Update)

	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoom().Delete)
}
