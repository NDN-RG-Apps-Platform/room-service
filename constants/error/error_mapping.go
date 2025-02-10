package error

import (
	errRoom "room-service/constants/error/room"
	errRoomSchedule "room-service/constants/error/roomSchedule"
)

func ErrMapping(err error) bool {
	allErrors := append(append(GeneralErrors, errRoom.RoomErrors...), errRoomSchedule.RoomScheduleErrors...)
	for _, item := range allErrors {
		if err == item {
			return true
		}
	}
	return false
}
