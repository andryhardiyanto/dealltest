package http

import (
	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	"github.com/AndryHardiyanto/dealltest/internal/transport/http/handler"
	libGin "github.com/AndryHardiyanto/dealltest/lib/gin"
	"github.com/gin-gonic/gin"
)

func Router(engineGin *gin.Engine, app *app.App) {
	engineGin.NoRoute(libGin.CustomPageNotFound())

	engineGin.Use(libGin.RegisterGinLogger())
	engineGin.Use(libGin.RegisterErrorHandler())
	engineGin.Use(libGin.Recovery())

	v1 := engineGin.Group("/v1")
	{
		v1.POST("/login", handler.Login(app))
		v1.POST("/renew/access", handler.RenewJwt(app))

		v1.Use(Middleware(app))
		v1.GET("/account", handler.Get(app))

		v1.Use(MiddlewareAdmin(app))
		v1.GET("/account/:id", handler.GetById(app))
		v1.POST("/account", handler.Register(app))
		v1.DELETE("/account/:id", handler.Delete(app))
		v1.PUT("/account/:id", handler.Update(app))

		v1.GET("/validate", handler.Validate(app))
	}

}
