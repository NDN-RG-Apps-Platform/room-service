package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type RoomRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Code        string                 `json:"code" validate:"required"`
	Capacity    string                 `json:"capacity" validate:"required"`
	Description string                 `json:"description" validate:"required"`
	Image       []multipart.FileHeader `json:"image" validate:"required"`
}

type UpdateFieldRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Code        string                 `json:"code" validate:"required"`
	Capacity    string                 `json:"capacity" validate:"required"`
	Description string                 `json:"description" validate:"required"`
	Image       []multipart.FileHeader `json:"image"`
}

type RoomResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Capacity    string    `json:"capacity"`
	Description string    `json:"description"`
	Image       []string  `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RoomDetailResponse struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Capacity    string    `json:"capacity"`
	Description string    `json:"description"`
	Image       []string  `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RoomRequestParam struct {
	Page       int     `json:"page" validates:"required"`
	Limit      int     `json:"limit" validates:"required"`
	SortColumn *string `json:"sortColumn"`
	SortOrder  *string `json:"sortOrder"`
	SetOrder   *string `json:"setOrder"`
}
