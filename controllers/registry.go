package controllers

import (
	controllers "room-service/controllers/room"
	controllers2 "room-service/controllers/roomSchedule"
	controllers3 "room-service/controllers/time"
	"room-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetRoom() controllers.IRoomController
	GetRoomSchedule() controllers2.IRoomScheduleController
	GetTime() controllers3.ITimeController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetRoom() controllers.IRoomController {
	return controllers.NewRoomController(r.service)
}

func (r *Registry) GetRoomSchedule() controllers2.IRoomScheduleController {
	return controllers2.NewRoomScheduleController(r.service)
}

func (r *Registry) GetTime() controllers3.ITimeController {
	return controllers3.NewTimeController(r.service)
}
