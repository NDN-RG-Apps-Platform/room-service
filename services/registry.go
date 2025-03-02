package services

import (
	"room-service/common/gcs"
	"room-service/repositories"
	roomService "room-service/services/room"
	roomScheduleService "room-service/services/roomSchedule"
	timeService "room-service/services/time"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IServiceRegistry interface {
	GetRoom() roomService.IRoomService
	GetRoomSchedule() roomScheduleService.IRoomScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &Registry{repository: repository, gcs: gcs}
}

func (r *Registry) GetRoom() roomService.IRoomService {
	return roomService.NewRoomService(r.repository, r.gcs)
}

func (r *Registry) GetRoomSchedule() roomScheduleService.IRoomScheduleService {
	return roomScheduleService.NewRoomScheduleService(r.repository)
}

func (r *Registry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(r.repository)
}
