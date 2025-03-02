package controllers

import (
	"net/http"
	"room-service/common/response"
	"room-service/domain/dto"
	"room-service/services"

	"github.com/gin-gonic/gin"
)

type TimeController struct {
	service services.IServiceRegistry
}

type ITimeController interface {
	GetAll(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
}

func NewTimeController(service services.IServiceRegistry) ITimeController {
	return &TimeController{service: service}
}

func (t *TimeController) GetAll(c *gin.Context) {
	result, err := t.service.GetTime().GetAll(c)
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

func (t *TimeController) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	result, err := t.service.GetTime().GetByUUID(c, uuid)
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

func (t *TimeController) Create(c *gin.Context) {
	var request dto.TimeRequest
	result, err := t.service.GetTime().Create(c, &request)
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
