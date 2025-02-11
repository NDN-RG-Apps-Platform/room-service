package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "room-service/common/error"
	errConstant "room-service/constants/error"
	errRoom "room-service/constants/error/room"
	"room-service/domain/dto"
	"room-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

type IRoomRepository interface {
	FindAllWithPagination(context.Context, *dto.RoomRequestParam) ([]models.Room, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.Room, error)
	FindByUUID(context.Context, string) (*models.Room, error)
	Create(context.Context, *models.Room) (*models.Room, error)
	Update(context.Context, string, *models.Room) (*models.Room, error)
	Delete(context.Context, string) error
}

func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &RoomRepository{db: db}
}

func (f *RoomRepository) FindAllWithPagination(ctx context.Context, param *dto.RoomRequestParam) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		sort  string
		total int64
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
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&rooms).
		Error

	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	err = f.db.
		WithContext(ctx).
		Model(&rooms).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return rooms, total, nil
}

func (f *RoomRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Room, error) {
	var rooms []models.Room
	err := f.db.
		WithContext(ctx).
		Find(&rooms).
		Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return rooms, nil
}

func (f *RoomRepository) FindByUUID(ctx context.Context, uuid string) (*models.Room, error) {
	var room models.Room
	err := f.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&room).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errRoom.ErrRoomNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &room, nil
}

func (f *RoomRepository) Create(ctx context.Context, req *models.Room) (*models.Room, error) {
	room := models.Room{
		UUID:        uuid.New(),
		Code:        req.Code,
		Name:        req.Name,
		Capacity:    req.Capacity,
		Description: req.Description,
		Image:       req.Image,
	}

	err := f.db.WithContext(ctx).Create(&room).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &room, nil
}

func (f *RoomRepository) Update(ctx context.Context, uuid string, req *models.Room) (*models.Room, error) {
	room := models.Room{
		Code:        req.Code,
		Name:        req.Name,
		Capacity:    req.Capacity,
		Description: req.Description,
		Image:       req.Image,
	}

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&room).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &room, nil
}

func (f *RoomRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Room{}).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
