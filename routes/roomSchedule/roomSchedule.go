package routes

import (
	"room-service/clients"
	"room-service/constants"
	"room-service/controllers"
	"room-service/middlewares"

	"github.com/gin-gonic/gin"
)

type RoomScheduleRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRoomScheduleRoute interface {
	Run()
}

func NewRoomScheduleRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRoomScheduleRoute {
	return &RoomScheduleRoute{controller: controller, group: group, client: client}
}

func (r *RoomScheduleRoute) Run() {
	group := r.group.Group("/room/schedule")
	group.GET("", middlewares.AuthenticateWithoutToken(), r.controller.GetRoomSchedule().GetAllByRoomIDAndDate)
	group.PATCH("", middlewares.AuthenticateWithoutToken(), r.controller.GetRoomSchedule().UpdateStatus)
	group.Use(middlewares.Authenticate())
	group.GET("/pagination", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
		constants.Lecture,
		constants.Student,
	}, r.client),
		r.controller.GetRoomSchedule().GetAllWithPagination)

	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
		constants.Lecture,
		constants.Student,
	}, r.client),
		r.controller.GetRoomSchedule().GetByUUID)

	group.POST("", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoomSchedule().Create)

	group.POST("/one-month", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoomSchedule().GenerateScheduleForOneMonth)

	group.PUT("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoomSchedule().Update)

	group.DELETE("/:uuid", middlewares.CheckRole([]string{
		constants.Administrator,
		constants.Co_Administrator,
		constants.Staff,
	}, r.client),
		r.controller.GetRoomSchedule().Delete)
}
