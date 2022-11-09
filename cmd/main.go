package main

import (
	"github.com/AndryHardiyanto/dealltest/config"
	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	repoAccount "github.com/AndryHardiyanto/dealltest/internal/repository/account"
	"github.com/AndryHardiyanto/dealltest/internal/service/account"
	"github.com/AndryHardiyanto/dealltest/internal/service/auth"
	transportHttp "github.com/AndryHardiyanto/dealltest/internal/transport/http"
	libLog "github.com/AndryHardiyanto/dealltest/lib/log"
	libPostgres "github.com/AndryHardiyanto/dealltest/lib/postgres"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.RegisterConfig()
	libLog.RegisterLogger(config.Cfg.Logger.Debug)
	gin.SetMode(config.Cfg.Server.GinMode)

	db := libPostgres.NewPostgres(config.Cfg.Database.DatabaseConnection)

	authService := auth.NewService(
		config.Cfg.Jwt.SignedSecret,
		config.Cfg.Jwt.AccessExpDuration,
		config.Cfg.Jwt.RefreshExpDuration,
	)
	accountRepo := repoAccount.NewPostgres(db)
	app := &app.App{
		Services: &app.Services{
			AuthService: authService,
			AccountService: account.NewService(
				accountRepo,
				authService,
			),
		},
	}

	transportHttp.RunServer(app,
		config.Cfg.Server.Port,
	)

}
