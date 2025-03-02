package controllers

import (
	"net/http"
	errValidation "room-service/common/error"
	"room-service/common/response"
	"room-service/domain/dto"
	"room-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type roomScheduleController struct {
	service services.IServiceRegistry
}

type IRoomScheduleController interface {
	GetAllWithPagination(c *gin.Context)
	GetAllByRoomIDAndDate(c *gin.Context)
	GetByUUID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	UpdateStatus(c *gin.Context)
	Delete(c *gin.Context)
	GenerateScheduleForOneMonth(c *gin.Context)
}

func NewRoomScheduleController(service services.IServiceRegistry) IRoomScheduleController {
	return &roomScheduleController{service: service}
}

func (f *roomScheduleController) GetAllWithPagination(c *gin.Context) {
	var params dto.RoomScheduleRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetRoomSchedule().GetAllWithPagination(c, &params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *roomScheduleController) GetAllByRoomIDAndDate(c *gin.Context) {
	var params dto.RoomScheduleByRoomIDAndDateRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetRoomSchedule().GetAllByRoomIDAndDate(c, c.Param("uuid"), params.Date)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *roomScheduleController) GetByUUID(c *gin.Context) {
	result, err := f.service.GetRoomSchedule().GetByUUID(c, c.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *roomScheduleController) Create(c *gin.Context) {
	var params dto.RoomScheduleRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	err = f.service.GetRoomSchedule().Create(c, &params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *roomScheduleController) GenerateScheduleForOneMonth(c *gin.Context) {
	var params dto.GenerateRoomScheduleForOneMostRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	err = f.service.GetRoomSchedule().GenerateScheduleForOneMonth(c, &params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *roomScheduleController) Update(c *gin.Context) {
	var params dto.UpdateRoomScheduleRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetRoomSchedule().Update(c, c.Param("uuid"), &params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Gin:  c,
		Data: result,
	})
}

func (f *roomScheduleController) UpdateStatus(c *gin.Context) {
	var request dto.UpdateStatusRoomScheduleRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})
		return
	}

	err = f.service.GetRoomSchedule().UpdateStatus(c, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *roomScheduleController) Delete(c *gin.Context) {
	err := f.service.GetRoomSchedule().Delete(c, c.Param("uuid"))
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Gin:  c,
	})
}
