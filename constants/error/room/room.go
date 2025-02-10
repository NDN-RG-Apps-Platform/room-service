package error

import "errors"

var (
	ErrRoomNotFound = errors.New("room not found")
)

var RoomErrors = []error{
	ErrRoomNotFound,
}
