package response

import (
	"net/http"
	"room-service/constants"
	errConstant "room-service/constants/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Token   *string     `json:"token,omitempty"`
}

type ParamHTTPResp struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    interface{}
	Token   *string
}

// Fungsi untuk mengirim response ke client
func HTTPResponse(param ParamHTTPResp) {
	// Jika tidak ada error, kirim response sukses
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})
		return
	}

	// Menentukan pesan error
	message := errConstant.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Err != nil {
		if errConstant.ErrMapping(param.Err) {
			message = param.Err.Error()
		}
	}

	// Kirim response error
	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    nil, // Data dikosongkan jika terjadi error untuk menjaga konsistensi response
	})
}
