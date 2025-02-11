package repositories

import (
	roomRepo "room-service/repositories/room"
	roomScheduleRepo "room-service/repositories/roomSchedule"
	timeRepo "room-service/repositories/time"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetRoom() roomRepo.IRoomRepository
	GetRoomSchedule() roomScheduleRepo.IRoomScheduleRepository
	GetTime() timeRepo.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetRoom() roomRepo.IRoomRepository {
	return roomRepo.NewRoomRepository(r.db)
}

func (r *Registry) GetRoomSchedule() roomScheduleRepo.IRoomScheduleRepository {
	return roomScheduleRepo.NewRoomScheduleRepository(r.db)
}

func (r *Registry) GetTime() timeRepo.ITimeRepository {
	return timeRepo.NewTimeRepository(r.db)
}
