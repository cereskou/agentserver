package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/db"
	"ditto.co.jp/agentserver/logger"
	"ditto.co.jp/agentserver/services"
	"ditto.co.jp/agentserver/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//
var (
	ErrFailedCreateMemDB = errors.New("failed to create memory db")
)

//RunMasterServer -
func RunMasterServer(conf *config.Config, db *db.Database) error {
	db.Map = util.NewAgentMap()
	if db.Map == nil {
		return ErrFailedCreateMemDB
	}

	if err := services.InitService("master", conf, db); err != nil {
		return err
	}
	defer services.Close()

	// Echo instance
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	//Log
	level := logger.GetLevel(conf.Log)
	e.Logger.SetLevel(log.Lvl(level))

	// Middleware
	if level == 0 {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Index: "index.html",
	}))

	//Routes
	services.MasterService.RegisterRoutes(e, "")

	//Run the server
	go func() {
		addr := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
		logger.Debugf("http server started on %v", addr)

		e.Logger.Fatal(e.Start(addr))
	}()

	//Start listen
	go util.RunListener(conf, db.Map)
	//Start Broadcast
	go util.RunBroadcast("master", conf)
	//Health check
	go util.HealthCheck(conf, db.Map)

	//Quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Printf("\r- Ctrl+C pressed in Terminal\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug("Close")
	//Broadcase
	//DB
	if err := db.Close(); err != nil {
		logger.Error(err)
	}
	//Echo
	if err := e.Close(); err != nil {
		logger.Error(err)
	}

	logger.Debug("Shutdown")
	if err := e.Shutdown(ctx); err != nil {
		logger.Error(err)
	}

	return nil
}
