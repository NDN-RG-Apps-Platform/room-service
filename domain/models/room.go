package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Room struct {
	ID   uint      `gorm:"primaryKey;autoIncrement"`
	UUID uuid.UUID `gorm:"type:uuid;not null"`
	// LibraryID     uint      `gorm:"type:int;not null"`
	Image         pq.StringArray `gorm:"type:text[];not null"`
	Code          string         `gorm:"type:varchar(15);not null"`
	Name          string         `gorm:"type:varchar(100);not null"`
	Capacity      string         `gorm:"type:varchar(15);not null"`
	Description   string         `gorm:"type:varchar(100);not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *gorm.DeletedAt
	RoomSchedules []RoomSchedule `gorm:"foreignKey:room_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
