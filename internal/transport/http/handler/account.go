package handler

import (
	"strconv"

	modelAccount "github.com/AndryHardiyanto/dealltest/internal/model/account"
	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	libGin "github.com/AndryHardiyanto/dealltest/lib/gin"
	"github.com/AndryHardiyanto/dealltest/lib/response"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

func Register(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *modelAccount.RegisterRequest

		err := libGin.NewRequest(c).GetBody(&req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetBody Register"))
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

		err = app.Services.AccountService.Register(c.Request.Context(), req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Register"))
			c.Abort()
			return
		}

		response.New(c).Created()
	}
}
func Login(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *modelAccount.LoginRequest

		err := libGin.NewRequest(c).GetBody(&req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetBody Login"))
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

		resp, err := app.Services.AccountService.Login(c.Request.Context(), req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Login"))
			c.Abort()
			return
		}

		response.New(c).Ok(resp)
	}
}
func Get(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accountID = c.GetInt64("account_id")

		resp, err := app.Services.AccountService.Get(c.Request.Context(), accountID)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Get"))
			c.Abort()
			return
		}

		response.New(c).Ok(resp)
	}
}
func GetById(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := libGin.NewRequest(c).GetParams()
		accountID, err := strconv.ParseInt(params.ByName("id"), 10, 64)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error ParseInt").SetType(errors.TypeUnprocessableEntity))
			c.Abort()
			return
		}

		resp, err := app.Services.AccountService.Get(c.Request.Context(), accountID)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Get"))
			c.Abort()
			return
		}

		response.New(c).Ok(resp)
	}
}
func Delete(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := libGin.NewRequest(c).GetParams()
		accountID, err := strconv.ParseInt(params.ByName("id"), 10, 64)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error ParseInt").SetType(errors.TypeUnprocessableEntity))
			c.Abort()
			return
		}

		err = app.Services.AccountService.Delete(c.Request.Context(), &modelAccount.DeleteRequest{
			AccountID: accountID,
		}, c.GetInt64("account_id"))
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Delete"))
			c.Abort()
			return
		}

		response.New(c).NonContent()
	}
}
func Update(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := libGin.NewRequest(c)
		params := request.GetParams()
		accountID, err := strconv.ParseInt(params.ByName("id"), 10, 64)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error ParseInt").SetType(errors.TypeUnprocessableEntity))
			c.Abort()
			return
		}
		var (
			req *modelAccount.UpdateRequest
		)
		err = request.GetBody(&req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetBody Update"))
			c.Abort()
			return
		}

		if req == nil {
			req = &modelAccount.UpdateRequest{}
		}
		req.AccountID = accountID

		err = app.Services.AccountService.Update(c.Request.Context(), req)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error Update"))
			c.Abort()
			return
		}

		response.New(c).NonContent()
	}
}
