package routes

import (
	"room-service/clients"
	"room-service/controllers"
	routes "room-service/routes/room"
	routes2 "room-service/routes/roomSchedule"
	timeRoute "room-service/routes/time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
	client     clients.IClientRegistry
}

type IRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup, client clients.IClientRegistry) IRegistry {
	return &Registry{controller: controller, group: group, client: client}
}

func (r *Registry) roomRoute() routes.IRoomRoute {
	return routes.NewRoomRoute(r.controller, r.group, r.client)
}

func (r *Registry) roomScheduleRoute() routes2.IRoomScheduleRoute {
	return routes2.NewRoomScheduleRoute(r.controller, r.group, r.client)
}

func (r *Registry) timeRoute() timeRoute.ITimeRoute {
	return timeRoute.NewTimeRoute(r.controller, r.group, r.client)
}

func (r *Registry) Serve() {
	r.roomRoute().Run()
	r.roomScheduleRoute().Run()
	r.timeRoute().Run()
}
