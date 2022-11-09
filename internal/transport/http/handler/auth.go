package handler

import (
	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	modelAuth "github.com/AndryHardiyanto/dealltest/internal/model/auth"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	libGin "github.com/AndryHardiyanto/dealltest/lib/gin"
	"github.com/AndryHardiyanto/dealltest/lib/response"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

func Validate(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.New(c).Ok(nil)
	}
}
func RenewJwt(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *modelAuth.RenewJwtRequest

		err := libGin.NewRequest(c).GetBody(&req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetBody RenewJwt"))
			c.Abort()
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error validate struct").SetType(errors.TypeRequestCannotEmpty))
			c.Abort()
			return
		}

		resp, err := app.Services.AuthService.RenewJwt(c.Request.Context(), req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error RenewJwt"))
			c.Abort()
			return
		}

		response.New(c).Ok(resp)
	}
}
