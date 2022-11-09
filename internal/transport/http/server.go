package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AndryHardiyanto/dealltest/internal/model/app"
	libLog "github.com/AndryHardiyanto/dealltest/lib/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServer(app *app.App, port string) {
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Api-Key", "X-CSRF-Token", "X-Access-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	Router(engine, app)

	libLog.Info().Msg(fmt.Sprintf("About to listen on %s. Go to http://127.0.0.1:%s", port, port))
	srv := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: engine}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			libLog.Error().Msg(fmt.Sprintf("err Listen: %s\n", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	libLog.Info().Msg("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err.Error())
		libLog.Panic().Msg(fmt.Sprintf("Server Shutdown: %s\n", err.Error()))
	}
	libLog.Info().Msg("Server exiting ...")
}
