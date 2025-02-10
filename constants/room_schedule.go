package constants

type RoomScheduleStatusName string
type RoomScheduleStatus int

const (
	Available RoomScheduleStatus = 100
	Booked    RoomScheduleStatus = 200

	AvailableString RoomScheduleStatusName = "Available"
	BookedString    RoomScheduleStatusName = "Booked"
)

var mapRoomScheduleStatusIntToString = map[RoomScheduleStatus]RoomScheduleStatusName{
	Available: AvailableString,
	Booked:    BookedString,
}

var mapRoomScheduleStatusStringToInt = map[RoomScheduleStatusName]RoomScheduleStatus{
	AvailableString: Available,
	BookedString:    Booked,
}

func (r RoomScheduleStatus) GetStatusString() RoomScheduleStatusName {
	return mapRoomScheduleStatusIntToString[r]
}

func (r RoomScheduleStatusName) GetStatusInt() RoomScheduleStatus {
	return mapRoomScheduleStatusStringToInt[r]
}
