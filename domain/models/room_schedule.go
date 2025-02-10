package models

import (
	"room-service/constants"
	"time"

	"github.com/google/uuid"
)

type RoomSchedule struct {
	ID        uint                         `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID                    `gorm:"type:uuid;not null"`
	RoomID    uint                         `gorm:"type:int;not null"`
	TimeID    uint                         `gorm:"type:int;not null"`
	Date      time.Time                    `gorm:"type:date;not null"`
	Status    constants.RoomScheduleStatus `gorm:"type:int;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	Room Room `gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Time Time `gorm:"foreignKey:time_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
