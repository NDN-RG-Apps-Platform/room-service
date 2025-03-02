package services

import (
	"context"
	"fmt"
	"room-service/common/util"
	"room-service/constants"
	errRoomSchedule "room-service/constants/error/roomSchedule"
	"room-service/domain/dto"
	"room-service/domain/models"
	"room-service/repositories"
	"time"

	"github.com/google/uuid"
)

type RoomScheduleService struct {
	repository repositories.IRepositoryRegistry
}

type IRoomScheduleService interface {
	GetAllWithPagination(context.Context, *dto.RoomScheduleRequestParam) (*util.PaginationResult, error)
	GetAllByRoomIDAndDate(context.Context, string, string) ([]dto.RoomScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.RoomScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateRoomScheduleForOneMostRequest) error
	Create(context.Context, *dto.RoomScheduleRequest) error
	Update(context.Context, string, *dto.UpdateRoomScheduleRequest) (*dto.RoomScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusRoomScheduleRequest) error
	Delete(context.Context, string) error
}

func NewRoomScheduleService(repository repositories.IRepositoryRegistry) IRoomScheduleService {
	return &RoomScheduleService{
		repository: repository,
	}
}

func (r *RoomScheduleService) GetAllWithPagination(
	ctx context.Context,
	param *dto.RoomScheduleRequestParam,
) (*util.PaginationResult, error) {
	roomSchedules, total, err := r.repository.GetRoomSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	roomScheduleResults := make([]dto.RoomScheduleResponse, 0, len(roomSchedules))
	for _, schedule := range roomSchedules {
		roomScheduleResults = append(roomScheduleResults, dto.RoomScheduleResponse{
			UUID:        schedule.UUID,
			RoomName:    schedule.Room.Name,
			Date:        schedule.Date.Format("2006-01-02"),
			Capacity:    schedule.Room.Capacity,
			Description: schedule.Room.Description,
			Status:      schedule.Status.GetStatusString(),
			Time:        fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:   *schedule.CreatedAt,
			UpdatedAt:   *schedule.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  roomScheduleResults,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (r *RoomScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	indonesiaMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Agu",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formattedDate := date.Format("02 Jan")
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indonesiaMonth[month])
	return formattedDate
}

func (r *RoomScheduleService) GetAllByRoomIDAndDate(ctx context.Context, uuid, date string) ([]dto.RoomScheduleForBookingResponse, error) {
	room, err := r.repository.GetRoom().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	roomSchedules, err := r.repository.GetRoomSchedule().FindAllByRoomIDAndDate(ctx, int(room.ID), date)
	if err != nil {
		return nil, err
	}

	roomScheduleResults := make([]dto.RoomScheduleForBookingResponse, 0, len(roomSchedules))
	for _, schedule := range roomSchedules {
		roomScheduleResults = append(roomScheduleResults, dto.RoomScheduleForBookingResponse{
			UUID:        schedule.UUID,
			Date:        r.convertMonthName(schedule.Date.Format(time.DateOnly)),
			Time:        schedule.Time.StartTime,
			Status:      schedule.Status.GetStatusString(),
			Capacity:    schedule.Room.Capacity,
			Description: schedule.Room.Description,
		})
	}

	return roomScheduleResults, nil
}

func (r *RoomScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.RoomScheduleResponse, error) {
	roomSchedule, err := r.repository.GetRoomSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.RoomScheduleResponse{
		UUID:        roomSchedule.UUID,
		RoomName:    roomSchedule.Room.Name,
		Date:        roomSchedule.Date.Format("2006-01-02"),
		Capacity:    roomSchedule.Room.Capacity,
		Description: roomSchedule.Room.Description,
		Status:      roomSchedule.Status.GetStatusString(),
		CreatedAt:   *roomSchedule.CreatedAt,
		UpdatedAt:   *roomSchedule.UpdatedAt,
	}

	return &response, nil
}

func (r *RoomScheduleService) Create(ctx context.Context, request *dto.RoomScheduleRequest) error {
	room, err := r.repository.GetRoom().FindByUUID(ctx, request.RoomID)
	if err != nil {
		return err
	}

	roomSchedules := make([]models.RoomSchedule, 0, len(request.TimeIDs))
	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	for _, timeID := range request.TimeIDs {
		scheduleTime, err := r.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		schedule, err := r.repository.GetRoomSchedule().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(room.ID))
		if err != nil {
			return err
		}

		if schedule != nil {
			return errRoomSchedule.ErrRoomScheduleIsExist
		}

		roomSchedules = append(roomSchedules, models.RoomSchedule{
			UUID:   uuid.New(),
			RoomID: room.ID,
			TimeID: scheduleTime.ID,
			Date:   dateParsed,
			Status: constants.Available,
		})
	}

	err = r.repository.GetRoomSchedule().Create(ctx, roomSchedules)
	if err != nil {
		return err
	}

	return nil

}

func (r *RoomScheduleService) GenerateScheduleForOneMonth(ctx context.Context, request *dto.GenerateRoomScheduleForOneMostRequest) error {
	room, err := r.repository.GetRoom().FindByUUID(ctx, request.RoomID)
	if err != nil {
		return err
	}

	timeSlots, err := r.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDays := 30
	roomSchedules := make([]models.RoomSchedule, 0, numberOfDays)
	now := time.Now().Add(24 * time.Hour)

	for i := 0; i < numberOfDays; i++ {
		currentDate := now.AddDate(0, 0, i)
		for _, item := range timeSlots {
			schedule, err := r.repository.GetRoomSchedule().FindByDateAndTimeID(ctx, currentDate.Format(time.DateOnly), int(item.ID), int(room.ID))
			if err != nil {
				return err
			}

			if schedule != nil {
				return errRoomSchedule.ErrRoomScheduleIsExist
			}

			roomSchedules = append(roomSchedules, models.RoomSchedule{
				UUID:   uuid.New(),
				RoomID: room.ID,
				TimeID: item.ID,
				Date:   currentDate,
				Status: constants.Available,
			})
		}
	}

	err = r.repository.GetRoomSchedule().Create(ctx, roomSchedules)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomScheduleService) Update(ctx context.Context, uuid string, request *dto.UpdateRoomScheduleRequest) (*dto.RoomScheduleResponse, error) {
	roomSchedule, err := r.repository.GetRoomSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	scheduleTime, err := r.repository.GetTime().FindByUUID(ctx, request.TimeID)
	if err != nil {
		return nil, err
	}

	isTimeExist, err := r.repository.GetRoomSchedule().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(roomSchedule.RoomID))
	if err != nil {
		return nil, err
	}

	if isTimeExist != nil && request.Date != roomSchedule.Date.Format(time.DateOnly) {
		checkDate, err := r.repository.GetRoomSchedule().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(roomSchedule.RoomID))
		if err != nil {
			return nil, err
		}

		if checkDate != nil {
			return nil, errRoomSchedule.ErrRoomScheduleIsExist
		}
	}

	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	roomResult, err := r.repository.GetRoomSchedule().Update(ctx, uuid, &models.RoomSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	response := dto.RoomScheduleResponse{
		UUID:        roomResult.UUID,
		RoomName:    roomResult.Room.Name,
		Date:        roomResult.Date.Format(time.DateOnly),
		Capacity:    roomResult.Room.Capacity,
		Description: roomResult.Room.Description,
		Status:      roomResult.Status.GetStatusString(),
		Time:        fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:   *roomResult.CreatedAt,
		UpdatedAt:   *roomResult.UpdatedAt,
	}
	return &response, nil
}

func (r *RoomScheduleService) UpdateStatus(ctx context.Context, request *dto.UpdateStatusRoomScheduleRequest) error {
	for _, item := range request.RoomScheduleIDs {
		_, err := r.repository.GetRoomSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = r.repository.GetRoomSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoomScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := r.repository.GetRoomSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = r.repository.GetRoomSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}
