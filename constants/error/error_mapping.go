package error

import (
	errRoom "room-service/constants/error/room"
	errRoomSchedule "room-service/constants/error/roomSchedule"
	errTime "room-service/constants/error/time"
)

func ErrMapping(err error) bool {
	var (
		GeneralErrors      = GeneralErrors
		RoomErrors         = errRoom.RoomErrors
		RoomScheduleErrors = errRoomSchedule.RoomScheduleErrors
		TimeErrors         = errTime.TimeErrors
	)

	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, RoomErrors...)
	allErrors = append(allErrors, RoomScheduleErrors...)
	allErrors = append(allErrors, TimeErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
