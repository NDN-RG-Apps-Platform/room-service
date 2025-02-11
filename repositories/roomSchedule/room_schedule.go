package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "room-service/common/error"
	constans "room-service/constants"
	errConstant "room-service/constants/error"
	errRoomSchedule "room-service/constants/error/roomSchedule"
	"room-service/domain/dto"
	"room-service/domain/models"

	"gorm.io/gorm"
)

type RoomScheduleRepository struct {
	db *gorm.DB
}

type IRoomScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.RoomScheduleRequestParam) ([]models.RoomSchedule, int64, error)
	FindAllByRoomIDAndDate(context.Context, int, string) ([]models.RoomSchedule, error)
	FindByUUID(context.Context, string) (*models.RoomSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.RoomSchedule, error)
	Create(context.Context, []models.RoomSchedule) error
	Update(context.Context, string, *models.RoomSchedule) (*models.RoomSchedule, error)
	UpdateStatus(context.Context, constans.RoomScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewRoomScheduleRepository(db *gorm.DB) IRoomScheduleRepository {
	return &RoomScheduleRepository{db: db}
}

func (f *RoomScheduleRepository) FindAllWithPagination(ctx context.Context, param *dto.RoomScheduleRequestParam) ([]models.RoomSchedule, int64, error) {
	var (
		roomSchedules []models.RoomSchedule
		sort          string
		total         int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.
		WithContext(ctx).
		Preload("Room").
		Preload("Time").
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&roomSchedules).
		Error

	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Model(&roomSchedules).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return roomSchedules, total, nil
}

func (f *RoomScheduleRepository) FindAllByRoomIDAndDate(ctx context.Context, roomID int, date string) ([]models.RoomSchedule, error) {
	var roomSchedules []models.RoomSchedule
	err := f.db.
		WithContext(ctx).
		Preload("Room").
		Preload("Time").
		Where("room_id = ?", roomID).
		Where("date = ?", date).
		Joins("LEFT JOIN times ON room_schedules.time_id = times.id").
		Order("times.start_time asc").
		Find(&roomSchedules).
		Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return roomSchedules, nil
}

func (f *RoomScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.RoomSchedule, error) {
	var roomSchedules models.RoomSchedule
	err := f.db.
		WithContext(ctx).
		Preload("Room").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&roomSchedules).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errRoomSchedule.ErrRoomScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &roomSchedules, nil
}

func (f *RoomScheduleRepository) FindByDateAndTimeID(ctx context.Context, date string, timeID int, roomID int) (*models.RoomSchedule, error) {
	var roomSchedules models.RoomSchedule
	err := f.db.
		WithContext(ctx).
		Where("date = ?", date).
		Where("time_id = ?", timeID).
		Where("room_id = ?", roomID).
		First(&roomSchedules).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errRoomSchedule.ErrRoomScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &roomSchedules, nil
}

func (f *RoomScheduleRepository) Create(ctx context.Context, req []models.RoomSchedule) error {
	err := f.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}

func (f *RoomScheduleRepository) Update(ctx context.Context, uuid string, req *models.RoomSchedule) (*models.RoomSchedule, error) {
	roomSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	roomSchedule.Date = req.Date
	err = f.db.WithContext(ctx).Save(&roomSchedule).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return roomSchedule, nil
}

func (f *RoomScheduleRepository) UpdateStatus(ctx context.Context, status constans.RoomScheduleStatus, uuid string) error {
	roomSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	roomSchedule.Status = status
	err = f.db.WithContext(ctx).Save(&roomSchedule).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	return nil
}

func (f *RoomScheduleRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.RoomSchedule{}).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
