package error

import "errors"

var (
	ErrRoomScheduleNotFound = errors.New("room schedule not found")
	ErrRoomScheduleIsExist  = errors.New("room schedule already exist")
)

var RoomScheduleErrors = []error{
	ErrRoomScheduleNotFound,
	ErrRoomScheduleIsExist,
}
