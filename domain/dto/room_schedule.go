package dto

import (
	"room-service/constants"
	"time"

	"github.com/google/uuid"
)

type RoomScheduleRequest struct {
	RoomID  string   `json:"RoomID" validate:"required"`
	Date    string   `json:"date" validate:"required"`
	TimeIDs []string `json:"timeIDs" validate:"required"`
}

type GenerateRoomScheduleForOneMostRequest struct {
	RoomID string `json:"RoomID" validate:"required"`
}

type UpdateRoomScheduleRequest struct {
	Date   string `json:"date" validate:"required"`
	TimeID string `json:"timeID" validate:"required"`
}

type UpdateStatusRoomScheduleRequest struct {
	RoomScheduleIDs []string `json:"roomScheduleIDs" validate:"required"`
}

type RoomScheduleResponse struct {
	UUID        uuid.UUID                        `json:"uuid"`
	RoomName    string                           `json:"roomName"`
	Capacity    string                           `json:"capacity"`
	Description string                           `json:"description"`
	Date        string                           `json:"date"`
	Status      constants.RoomScheduleStatusName `json:"status"`
	Time        string                           `json:"time"`
	CreatedAt   time.Time                        `json:"createdAt"`
	UpdatedAt   time.Time                        `json:"updatedAt"`
}

type RoomScheduleForBookingResponse struct {
	UUID        uuid.UUID                        `json:"uuid"`
	Date        string                           `json:"date"`
	Capacity    string                           `json:"capacity"`
	Description string                           `json:"description"`
	Status      constants.RoomScheduleStatusName `json:"status"`
	Time        string                           `json:"time"`
}

type RoomScheduleRequestParam struct {
	Page       int     `json:"page" validates:"required"`
	Limit      int     `json:"limit" validates:"required"`
	SortColumn *string `json:"sortColumn"`
	SortOrder  *string `json:"sortOrder"`
	SetOrder   *string `json:"setOrder"`
}

type RoomScheduleByRoomIDAndDateRequestParam struct {
	Date string `json:"date" validate:"required"`
}
