package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSONResponse struct {
	Code int
	Msg string
	Data interface{}
}

func (jr *JSONResponse) FailJSON(ct *gin.Context, errorMsg string, data interface{})  {
	resultData := JSONResponse{
		Code : 0,
		Msg : errorMsg,
		Data : data,
	}
	ct.JSON(http.StatusOK,resultData)
}

func (jr *JSONResponse) SuccJSON(ct *gin.Context, succMsg string, data interface{})  {
	resultData := JSONResponse{
		Code : 1,
		Msg : succMsg,
		Data : data,
	}
	ct.JSON(http.StatusOK,resultData)
}
