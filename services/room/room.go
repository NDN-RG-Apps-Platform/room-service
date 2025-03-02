package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"room-service/common/gcs"
	"room-service/common/util"
	errConstant "room-service/constants/error"
	"room-service/domain/dto"
	"room-service/domain/models"
	"room-service/repositories"
	"time"
)

type RoomService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
}

type IRoomService interface {
	GetAllWithPagination(context.Context, *dto.RoomRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(context.Context) ([]dto.RoomResponse, error)
	GetByUUID(context.Context, string) (*dto.RoomResponse, error)
	Create(context.Context, *dto.RoomRequest) (*dto.RoomResponse, error)
	Update(context.Context, string, *dto.RoomRequest) (*dto.RoomResponse, error)
	Delete(context.Context, string) error
}

func NewRoomService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient) IRoomService {
	return &RoomService{repository: repository, gcs: gcs}
}

func (r *RoomService) GetAllWithPagination(ctx context.Context, param *dto.RoomRequestParam) (*util.PaginationResult, error) {
	rooms, total, err := r.repository.GetRoom().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	roomResults := make([]dto.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomResults = append(roomResults, dto.RoomResponse{
			UUID:        room.UUID,
			Code:        room.Code,
			Name:        room.Name,
			Capacity:    room.Capacity,
			Description: room.Description,
			Image:       room.Image,
			CreatedAt:   *room.CreatedAt,
			UpdatedAt:   *room.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  roomResults,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (r *RoomService) GetAllWithoutPagination(ctx context.Context) ([]dto.RoomResponse, error) {
	rooms, err := r.repository.GetRoom().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	roomResults := make([]dto.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomResults = append(roomResults, dto.RoomResponse{
			UUID:        room.UUID,
			Name:        room.Name,
			Capacity:    room.Capacity,
			Description: room.Description,
			Image:       room.Image,
		})
	}

	return roomResults, nil
}

func (r *RoomService) GetByUUID(ctx context.Context, uuid string) (*dto.RoomResponse, error) {
	room, err := r.repository.GetRoom().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	roomResult := dto.RoomResponse{
		UUID:        room.UUID,
		Code:        room.Code,
		Name:        room.Name,
		Capacity:    room.Capacity,
		Description: room.Description,
		Image:       room.Image,
		CreatedAt:   *room.CreatedAt,
		UpdatedAt:   *room.UpdatedAt,
	}

	return &roomResult, nil
}

// upload file ke GCS
func (r *RoomService) validateUpload(images []multipart.FileHeader) error {
	if len(images) == 0 {
		return errConstant.ErrInvalidUploadFile
	}
	// maxima mb
	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return errConstant.ErrSizeTooBig
		}
	}

	return nil
}

func (r *RoomService) processAndUploadImage(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("Images%s-%s-%s", time.Now().Format("2006-01-02"), image.Filename, path.Ext(image.Filename))
	url, err := r.gcs.UploadFile(ctx, fileName, buffer.Bytes())
	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *RoomService) uploadImage(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := r.validateUpload(images)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := r.processAndUploadImage(ctx, image)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (r *RoomService) Create(ctx context.Context, request *dto.RoomRequest) (*dto.RoomResponse, error) {
	imageUrl, err := r.uploadImage(ctx, request.Image)
	if err != nil {
		return nil, err
	}

	room, err := r.repository.GetRoom().Create(ctx, &models.Room{
		Code:        request.Code,
		Name:        request.Name,
		Capacity:    request.Capacity,
		Description: request.Description,
		Image:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RoomResponse{
		UUID:        room.UUID,
		Code:        room.Code,
		Name:        room.Name,
		Capacity:    room.Capacity,
		Description: room.Description,
		Image:       room.Image,
		CreatedAt:   *room.CreatedAt,
		UpdatedAt:   *room.UpdatedAt,
	}
	return response, nil
}

func (r *RoomService) Update(ctx context.Context, uuid string, request *dto.RoomRequest) (*dto.RoomResponse, error) {
	room, err := r.repository.GetRoom().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var imageUrl []string
	if request.Image == nil {
		imageUrl = room.Image
	} else {
		imageUrl, err = r.uploadImage(ctx, request.Image)
		if err != nil {
			return nil, err
		}
	}

	roomResult, err := r.repository.GetRoom().Update(ctx, uuid, &models.Room{
		Code:        request.Code,
		Name:        request.Name,
		Capacity:    request.Capacity,
		Description: request.Description,
		Image:       imageUrl,
	})
	if err != nil {
		return nil, err
	}

	return &dto.RoomResponse{
		UUID:        roomResult.UUID,
		Code:        roomResult.Code,
		Name:        roomResult.Name,
		Capacity:    roomResult.Capacity,
		Description: roomResult.Description,
		Image:       roomResult.Image,
		CreatedAt:   *roomResult.CreatedAt,
		UpdatedAt:   *roomResult.UpdatedAt,
	}, nil
}

func (r *RoomService) Delete(ctx context.Context, uuid string) error {
	_, err := r.repository.GetRoom().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = r.repository.GetRoom().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
