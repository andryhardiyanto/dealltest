package response

import (
	"net/http"

	"github.com/AndryHardiyanto/dealltest/lib/errors"

	"github.com/gin-gonic/gin"
)

type response struct {
	Type    string       `json:"type"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
	context *gin.Context `json:"-"`
}

type Response interface {
	Ok(data interface{})
	Created()
	NonContent()
	Error(err error)
}

func New(c *gin.Context) Response {
	return &response{
		context: c,
	}
}

//Ok, For Ok data
func (r *response) Ok(data interface{}) {
	r.Data = data
	r.Message = "Success"
	r.Type = "OK"
	r.context.JSON(http.StatusOK, r)
}

//Created, For method POST
func (r *response) Created() {
	r.Message = "Success"
	r.Type = "CREATED"
	r.context.JSON(http.StatusCreated, r)
}

//NonContent, For method PUT,PATCH AND DELETE
func (r *response) NonContent() {
	r.Message = "Success"
	r.Type = "NON_CONTENT"
	r.context.JSON(http.StatusNoContent, r)
}

//Error, For method PUT,PATCH AND DELETE
func (r *response) Error(err error) {
	errorMessage := errors.RegisterErrorMessage()
	msg := errorMessage.Translate(r.context.Request, err)
	r.context.JSON(msg.Code, msg)
}
