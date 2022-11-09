package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	"github.com/AndryHardiyanto/dealltest/internal/model/auth"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/gin-gonic/gin"
)

func Middleware(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c.Request)
		if token == "" {
			c.Error(errors.NewError("validate empty token").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		err := app.Services.AuthService.Validate(c.Request.Context(), &auth.ValidateRequest{
			Token: token,
		})
		if err != nil {
			c.Error(errors.NewWrapError(err, "error validate token"))
			c.Abort()
			return
		}

		dataClaims, err := app.Services.AuthService.GetClaims(c.Request.Context(), &auth.GetClaimsRequest{
			Token: token,
		})
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetClaims"))
			c.Abort()
			return
		}

		id, err := strconv.ParseInt(dataClaims.Subject, 10, 64)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error ParseInt").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		if id <= 0 {
			c.Error(errors.NewWrapError(err, "error invalid account id").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		c.Set("account_id", id)

		c.Next()
	}
}

func MiddlewareAdmin(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c.Request)
		if token == "" {
			c.Error(errors.NewError("validate empty token").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		err := app.Services.AuthService.Validate(c.Request.Context(), &auth.ValidateRequest{
			Token: token,
		})
		if err != nil {
			c.Error(errors.NewWrapError(err, "error validate token"))
			c.Abort()
			return
		}

		dataClaims, err := app.Services.AuthService.GetClaims(c.Request.Context(), &auth.GetClaimsRequest{
			Token: token,
		})
		if err != nil {
			c.Error(errors.NewWrapError(err, "error GetClaims"))
			c.Abort()
			return
		}

		id, err := strconv.ParseInt(dataClaims.Subject, 10, 64)
		if err != nil {
			c.Error(errors.NewWrapError(err, "error ParseInt").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		if id <= 0 {
			c.Error(errors.NewWrapError(err, "error invalid account id").SetType(errors.TypeUnauthorized))
			c.Abort()
			return
		}

		if dataClaims.Role != "admin" {
			c.Error(errors.NewError("error permission denied").SetType(errors.TypePermissionDenied))
			c.Abort()
			return
		}

		c.Set("account_id", id)
		c.Set("role", dataClaims.Role)
		c.Next()
	}
}

// getToken .
func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) < 2 {
		return ""
	}

	token = strings.Trim(splitToken[1], " ")
	return token
}
